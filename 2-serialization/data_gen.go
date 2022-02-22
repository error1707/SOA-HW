package main

import "math/rand"

const bigText = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Cras nisl mauris, mattis eget felis quis, rutrum laoreet mauris. Mauris nec nibh nunc. Nullam facilisis, mauris id tempor tempor, lorem est eleifend velit, nec ultrices leo massa a velit. Curabitur a mauris ligula. Ut quis nulla nec enim ultricies consectetur. Morbi arcu nibh, rhoncus id nulla ut, fermentum mollis elit. Maecenas lacinia neque dolor, ac rhoncus leo tempus in. Fusce congue tortor eget viverra tempus. In fermentum elit vel augue ultrices pharetra.
Vestibulum ac lorem a sem aliquam scelerisque et at massa. Fusce vehicula erat nec metus gravida posuere. Praesent eu purus eget massa tincidunt feugiat. Cras rhoncus tincidunt massa iaculis scelerisque. Quisque ultricies metus quis neque fringilla pellentesque. Mauris ac cursus massa, eu ullamcorper odio. Nullam quis sapien sem. Vestibulum quis congue sapien, vel pulvinar arcu. Ut at metus quis felis laoreet mollis nec non velit. Nullam mattis elementum tortor eu commodo. Quisque consequat accumsan urna, non vehicula purus luctus a. Maecenas interdum finibus dictum.
Aenean sed sollicitudin elit. Vivamus facilisis eu libero ac ullamcorper. Aliquam erat volutpat. Donec in pulvinar dolor, sed iaculis sapien. Quisque vitae feugiat ligula. Aenean rutrum magna odio, ut viverra dui aliquet sit amet. Aenean in scelerisque justo. Suspendisse cursus ipsum at lorem vehicula, sed sodales est dictum.
Fusce congue eros sed lacus pellentesque, ac suscipit enim auctor. Nunc volutpat urna sed purus lobortis cursus. Donec ipsum est, condimentum eget elit sed, consectetur condimentum nunc. Morbi at sodales est. Donec rhoncus lectus lacus, vel molestie elit molestie in. Aenean id sem maximus, euismod dui sit amet, efficitur velit. Phasellus mollis non dolor a tristique. Ut vitae erat a neque tincidunt elementum. Nulla facilisi. Phasellus sit amet sem et tellus ornare congue id eget arcu.
Vivamus sed libero ut libero tempus pharetra at id nibh. Ut eget velit feugiat, facilisis dui vitae, vulputate libero. Morbi sed varius est, et cursus mauris. Aenean sagittis gravida ipsum, vitae consectetur libero sodales eget. Donec semper orci vel finibus scelerisque. Praesent eu eleifend eros, at egestas enim. Quisque placerat tincidunt finibus. Vivamus elementum consequat purus. Suspendisse a neque mi. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam ut felis ac eros rutrum rhoncus. Maecenas eu diam ut neque molestie rutrum. Nunc sit amet ante at nulla eleifend vehicula. Nulla in volutpat lacus. Integer euismod arcu velit, ac malesuada ligula vestibulum at.`

type StringMap map[string]string

type SmallStruct struct {
	SomeString string `xml:"some_string" yaml:"some-string"`
	SomeInt    int32  `xml:"some_int" yaml:"some-int"`
}

type BigStruct struct {
	SomeText  string        `xml:"some_text" yaml:"some-text"`
	SomeInt   int32         `xml:"some_int" yaml:"some-int"`
	SomeFloat float32       `xml:"some_float" yaml:"some-float"`
	SomeMap   StringMap     `xml:"some_map" yaml:"some-map"`
	SomeArray []SmallStruct `xml:"some_array" yaml:"some-array"`
}

func randString() string {
	b := make([]rune, 10)
	for i := range b {
		b[i] = rune('a' + rand.Intn(26))
	}
	return string(b)
}

func genMap(n int) StringMap {
	res := StringMap{}
	for ; n != 0; n-- {
		res[randString()] = randString()
	}
	return res
}

func genSmallStruct() SmallStruct {
	return SmallStruct{
		SomeString: randString(),
		SomeInt:    rand.Int31(),
	}
}

func genBigStruct() BigStruct {
	var arr []SmallStruct
	for i := 0; i < 10; i++ {
		arr = append(arr, genSmallStruct())
	}
	return BigStruct{
		SomeText:  bigText,
		SomeInt:   rand.Int31(),
		SomeFloat: rand.Float32(),
		SomeMap:   genMap(20),
		SomeArray: arr,
	}
}

var testBigStruct = genBigStruct()
