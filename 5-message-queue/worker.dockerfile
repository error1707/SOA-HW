FROM python:3.8-alpine

ADD worker ./

RUN pip3 install -r requirements.txt

CMD ["python3", "worker.py"]
#CMD ["/bin/sh"]