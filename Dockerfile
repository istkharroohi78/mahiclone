FROM golang:latest

# Ffmpeg install karna music streaming ke liye
RUN apt-get update && \
    apt-get install -y ffmpeg && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app/

# Saara code container me copy karna
COPY . /app/

# Naya Go version ab saare packages bina error ke download kar lega
RUN go mod tidy
RUN go build -o shivmusic_bot .

# Bot ko start karne ki command
CMD ["./shivmusic_bot"]
