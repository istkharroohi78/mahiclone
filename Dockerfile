FROM golang:1.22-bullseye

# Ffmpeg install karna music streaming ke liye
RUN apt-get update && \
    apt-get install -y ffmpeg && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app/

# Saara code container me copy karna
COPY . /app/

# Go ke packages download aur build karna
RUN go mod download
RUN go build -o shivmusic_bot .

# Bot ko start karne ki command
CMD ["./shivmusic_bot"]