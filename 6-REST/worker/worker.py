import json

import pika
from pdfdocument.document import PDFDocument
from uuid import uuid4


def generate_report(user):
    file_name = f"{uuid4()}.pdf"
    with open(f"./reports/{file_name}", "wb") as f:
        pdf = PDFDocument(f)
        pdf.init_report()
        pdf.h1('Report')
        pdf.p(f'Username: {user["username"]}')
        pdf.p(f'Email: {user["email"]}')
        pdf.p(f'Gender: {"Male" if user["gender"] == "m" else "Female"}')
        pdf.p(f'Win count: {user["win_count"]}')
        pdf.p(f'Lose count: {user["lose_count"]}')
        pdf.p(f'Time in game: {user["time_in_game"]}')
        pdf.generate()
    return file_name


def main():
    connection = pika.BlockingConnection(pika.ConnectionParameters('rabbitmq'))
    channel = connection.channel()
    channel.queue_declare(queue='requests')
    channel.queue_declare(queue='responses')

    def callback(ch, method, properties, body):
        req = json.loads(body)
        print(f"Received {req}", flush=True)
        resp = {
            "id": req["id"],
            "file_path": generate_report(req["user"])
        }
        channel.basic_publish(exchange='',
                              routing_key='responses',
                              body=json.dumps(resp).encode('UTF-8'))
    print("Ready", flush=True)
    channel.basic_consume(queue='requests',
                          auto_ack=True,
                          on_message_callback=callback)
    channel.start_consuming()


main()
