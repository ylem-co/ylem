FROM python:3.12-alpine

RUN apk add build-base libressl-dev libffi-dev python3-dev

WORKDIR /opt/ylem_python_processor
COPY . .
RUN pip install -r requirements.txt

EXPOSE 7338

CMD ["uvicorn", "--host", "0.0.0.0", "--port", "7338", "main:app"]
