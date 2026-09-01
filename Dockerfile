FROM golang:1.22-bullseye

# Ffmpeg install karna music streaming ke liye
RUN apt-get update && \
    apt-get install -y ffmpeg && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app/

# Saara code container me copy karna
COPY . /app/

# 🚀 FIX: go mod tidy saare missing dependencies aur go.sum errors ko auto-fix kar dega
RUN go mod tidy
RUN go build -o shivmusic_bot .

# Bot ko start karne ki command
CMD ["./shivmusic_bot"]
