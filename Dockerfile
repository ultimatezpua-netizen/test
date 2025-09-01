FROM python:3.11-slim

ENV PYTHONUNBUFFERED=1 \
    PIP_NO_CACHE_DIR=1

RUN useradd -m appuser
WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY . .

EXPOSE 8080

USER appuser

HEALTHCHECK CMD python -c "import socket; s=socket.socket(); s.settimeout(2); s.connect(('127.0.0.1', 8080)); s.close()" || exit 1

CMD ["python", "main.py"]