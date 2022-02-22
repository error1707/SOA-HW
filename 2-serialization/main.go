package main

import (
	pb "SOA/protobuf"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	avro "github.com/leboncoin/avrocado"
	"github.com/vmihailenco/msgpack/v5"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"reflect"
	"time"
)

func (m StringMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(m) == 0 {
		return nil
	}
	err := e.EncodeToken(start)
	if err != nil {
		return err
	}
	for k, v := range m {
		e.Encode(xmlMapEntry{XMLName: xml.Name{Local: k}, Value: v})
	}
	return e.EncodeToken(start.End())
}

func (m *StringMap) UnmarshalXML(d *xml.Decoder, _ xml.StartElement) error {
	*m = StringMap{}
	for {
		var e xmlMapEntry
		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		(*m)[e.XMLName.Local] = e.Value
	}
	return nil
}

type Encodable interface {
	Encode(interface{}) error
}

type Decodable interface {
	Decode(interface{}) error
}

func GetFile(name string) *os.File {
	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		panic(err)
	}
	return file
}

func MeasureBigStruct() (map[string]int64, map[string]int64) {
	var enc Encodable
	var dec Decodable
	var res BigStruct
	timeMeasurements := map[string]int64{}
	byteMeasurements := map[string]int64{}

	// -------------------- Native --------------------
	file := GetFile(`./objects/obj.gob`)
	enc = gob.NewEncoder(file)
	start := time.Now()
	err := enc.Encode(testBigStruct)
	timeMeasurements["NativeSerializeBigStruct"] = time.Since(start).Microseconds()
	if err != nil {
		panic(err)
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		panic(err)
	}

	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}
	byteMeasurements["NativeBigStruct"] = stat.Size()

	dec = gob.NewDecoder(file)
	start = time.Now()
	err = dec.Decode(&res)
	timeMeasurements["NativeDeserializeBigStruct"] = time.Since(start).Microseconds()
	if err != nil {
		panic(err)
	}
	err = file.Close()
	if err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(testBigStruct, res) {
		panic(`Native BigStruct not equal`)
	}

	// -------------------- XML --------------------
	res = BigStruct{}
	file = GetFile(`./objects/obj.xml`)
	enc = xml.NewEncoder(file)
	start = time.Now()
	err = enc.Encode(testBigStruct)
	timeMeasurements["XMLSerializeBigStruct"] = time.Since(start).Microseconds()
	if err != nil {
		panic(err)
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		panic(err)
	}

	stat, err = file.Stat()
	if err != nil {
		panic(err)
	}
	byteMeasurements["XMLBigStruct"] = stat.Size()

	dec = xml.NewDecoder(file)
	start = time.Now()
	err = dec.Decode(&res)
	timeMeasurements["XMLDeserializeBigStruct"] = time.Since(start).Microseconds()
	if err != nil {
		panic(err)
	}
	err = file.Close()
	if err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(testBigStruct, res) {
		panic(`XML BigStruct not equal`)
	}

	// -------------------- JSON --------------------
	res = BigStruct{}
	file = GetFile(`./objects/obj.json`)
	enc = json.NewEncoder(file)
	start = time.Now()
	err = enc.Encode(testBigStruct)
	timeMeasurements["JSONSerializeBigStruct"] = time.Since(start).Microseconds()
	if err != nil {
		panic(err)
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		panic(err)
	}

	stat, err = file.Stat()
	if err != nil {
		panic(err)
	}
	byteMeasurements["JSONBigStruct"] = stat.Size()

	dec = json.NewDecoder(file)
	start = time.Now()
	err = dec.Decode(&res)
	timeMeasurements["JSONDeserializeBigStruct"] = time.Since(start).Microseconds()
	if err != nil {
		panic(err)
	}
	err = file.Close()
	if err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(testBigStruct, res) {
		panic(`JSON BigStruct not equal`)
	}

	// -------------------- YAML --------------------
	res = BigStruct{}
	file = GetFile(`./objects/obj.yaml`)
	enc = yaml.NewEncoder(file)
	start = time.Now()
	err = enc.Encode(testBigStruct)
	timeMeasurements["YAMLSerializeBigStruct"] = time.Since(start).Microseconds()
	if err != nil {
		panic(err)
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		panic(err)
	}

	stat, err = file.Stat()
	if err != nil {
		panic(err)
	}
	byteMeasurements["YAMLBigStruct"] = stat.Size()

	dec = yaml.NewDecoder(file)
	start = time.Now()
	err = dec.Decode(&res)
	timeMeasurements["YAMLDeserializeBigStruct"] = time.Since(start).Microseconds()
	if err != nil {
		panic(err)
	}
	err = file.Close()
	if err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(testBigStruct, res) {
		panic(`YAML BigStruct not equal`)
	}

	// -------------------- Protobuf --------------------
	file = GetFile(`./objects/obj.pb`)
	var arr []*pb.SmallStruct
	for _, elem := range testBigStruct.SomeArray {
		arr = append(arr, &pb.SmallStruct{
			SomeInt:    elem.SomeInt,
			SomeString: elem.SomeString,
		})
	}
	protoStruct := pb.BigStruct{
		SomeText:  testBigStruct.SomeText,
		SomeInt:   testBigStruct.SomeInt,
		SomeFloat: testBigStruct.SomeFloat,
		SomeMap:   testBigStruct.SomeMap,
		SomeArray: arr,
	}
	start = time.Now()
	data, err := proto.Marshal(&protoStruct)
	timeMeasurements["ProtobufSerializeBigStruct"] = time.Since(start).Microseconds()
	if err != nil {
		panic(err)
	}
	_, err = file.Write(data)
	if err != nil {
		panic(err)
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		panic(err)
	}

	stat, err = file.Stat()
	if err != nil {
		panic(err)
	}
	byteMeasurements["ProtobufBigStruct"] = stat.Size()

	data = make([]byte, len(data))
	_, err = file.Read(data)
	if err != nil {
		panic(err)
	}
	protoRes := pb.BigStruct{}
	start = time.Now()
	err = proto.Unmarshal(data, &protoRes)
	timeMeasurements["ProtobufDeserializeBigStruct"] = time.Since(start).Microseconds()
	if err != nil {
		panic(err)
	}

	if !proto.Equal(&protoStruct, &protoRes) {
		panic(`Proto BigStruct not equal`)
	}

	// -------------------- MessagePack --------------------
	res = BigStruct{}
	file = GetFile(`./objects/obj.mp`)
	enc = msgpack.NewEncoder(file)
	start = time.Now()
	err = enc.Encode(testBigStruct)
	timeMeasurements["MessagePackSerializeBigStruct"] = time.Since(start).Microseconds()
	if err != nil {
		panic(err)
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		panic(err)
	}

	stat, err = file.Stat()
	if err != nil {
		panic(err)
	}
	byteMeasurements["MessagePackBigStruct"] = stat.Size()

	dec = msgpack.NewDecoder(file)
	start = time.Now()
	err = dec.Decode(&res)
	timeMeasurements["MessagePackDeserializeBigStruct"] = time.Since(start).Microseconds()
	if err != nil {
		panic(err)
	}
	err = file.Close()
	if err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(testBigStruct, res) {
		panic(`MessagePack BigStruct not equal`)
	}

	// -------------------- Avro --------------------
	res = BigStruct{}
	file = GetFile(`./objects/obj.avro`)
	codec, err := avro.NewCodec(`{
    	"type": "record",
    	"name": "BigStruct",
    	"fields": [
        	{"name": "SomeText", "type": "string"},
			{"name": "SomeInt", "type": "int"},
			{"name": "SomeFloat", "type": "float"},
			{"name": "SomeMap", "type": "map", "values": "string"},
			{"name": "SomeArray", "type": "array", "items": {
				"type": "record",
				"name": "SmallStruct",
				"fields": [
					{"name": "SomeString", "type": "string"},
					{"name": "SomeInt", "type": "int"}
				]
			}}
    	]
	}`)
	if err != nil {
		panic(err)
	}
	start = time.Now()
	data, err = codec.Marshal(testBigStruct)
	timeMeasurements["AvroSerializeBigStruct"] = time.Since(start).Microseconds()
	if err != nil {
		panic(err)
	}
	_, err = file.Write(data)
	if err != nil {
		panic(err)
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		panic(err)
	}

	stat, err = file.Stat()
	if err != nil {
		panic(err)
	}
	byteMeasurements["AvroBigStruct"] = stat.Size()

	data = make([]byte, len(data))
	_, err = file.Read(data)
	if err != nil {
		panic(err)
	}
	start = time.Now()
	err = codec.Unmarshal(data, &res)
	timeMeasurements["AvroDeserializeBigStruct"] = time.Since(start).Microseconds()
	if !reflect.DeepEqual(testBigStruct, res) {
		panic(`Avro BigStruct not equal`)
	}

	return timeMeasurements, byteMeasurements
}

func repeatFuncN(f func() (map[string]int64, map[string]int64), n int) (map[string]int64, map[string]int64) {
	tRes := map[string]int64{}
	var bRes, t map[string]int64
	for i := 0; i < n; i++ {
		t, bRes = f()
		for k, v := range t {
			tRes[k] = tRes[k] + v
		}
	}
	for k, v := range tRes {
		tRes[k] = v / int64(n)
	}
	return tRes, bRes
}

func SerializeResults(obj interface{}, name string) {
	file, err := os.OpenFile(`./results/`+name+`.json`, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		panic(err)
	}
	enc := json.NewEncoder(file)
	enc.SetIndent("", "    ")
	err = enc.Encode(obj)
	if err != nil {
		panic(err)
	}
}

type xmlMapEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

func main() {
	t, b := repeatFuncN(MeasureBigStruct, 10000)
	SerializeResults(t, "time_results")
	SerializeResults(b, "byte_results")
}
