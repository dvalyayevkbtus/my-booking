package payment

import (
	"dvalyayevkbtu/my-booking/config"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type InvoiceCreated struct {
	Reference string `json:"reference"`
	Volume    string `json:"volume"`
	Currency  string `json:"currency"`
}

type Invoice struct {
	Reference       string       `json:"reference"`
	Volume          string       `json:"volume"`
	Currency        string       `json:"currency"`
	VolumeFulfilled string       `json:"volumeFulfilled"`
	Status          string       `json:"status"`
	Confirments     []Confirment `json:"confirments"`
}

type Confirment struct {
	Reference   string `json:"reference"`
	Volume      string `json:"volume"`
	Currency    string `json:"currency"`
	AccountCode string `json:"accountCode"`
}

type Payment struct {
	baseUrl string
	client  *http.Client
}

func CreatePayment(conf config.PaymentConfig) *Payment {
	return &Payment{conf.URL, &http.Client{}}
}

func (p *Payment) CreateInvoice(reference, volume, currency string) error {
	reqStr, jsonErr := json.Marshal(InvoiceCreated{reference, volume, currency})
	if jsonErr != nil {
		return jsonErr
	}

	url := fmt.Sprintf("%s/payment", p.baseUrl)
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(reqStr)))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return errors.New("status is not accepted")
	}
	return nil
}

func (p *Payment) CheckPayment(reference string) (bool, error) {
	url := fmt.Sprintf("%s/payment/%s", p.baseUrl, reference)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false, errors.New("unsuccessful status code of check payment response")
	}

	str, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	var invoice Invoice
	err = json.Unmarshal(str, &invoice)
	if err != nil {
		return false, err
	}
	return invoice.Status == "FULFILLED", nil
}



TASK 1
bashsudo groupadd sysadm
sudo groupadd operations

sudo vi /etc/sudoers.d/sysadm

%sysadm ALL=(ALL) NOPASSWD: /bin/systemctl

sudo chmod 440 /etc/sudoers.d/sysadm
sudo cp /etc/sudoers.d/sysadm /home/kbtu/sysadm.bak

sudo vi /etc/sudoers.d/operations

%operations ALL=(ALL) NOPASSWD: /usr/local/bin/build-pipeline.sh
sudo chmod 440 /etc/sudoers.d/operations
sudo cp /etc/sudoers.d/operations /home/kbtu/operations.bak

TASK 2
sudo dnf install -y git

git clone https://github.com/dvalyayevkbtu/payment ~/payment
cd ~/payment
cat README.md

sudo vi ~/payment/Containerfile

FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production
COPY . .
RUN npm run build --if-present

FROM node:20-alpine
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN find / -perm /6000 -type f -exec chmod a-s {} \; 2>/dev/null || true
WORKDIR /app
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/node_modules ./node_modules
USER appuser
EXPOSE 4200
CMD ["node", "dist/server.js"]

sudo cp ~/payment/Containerfile /home/kbtu/Containerfile

sudo podman build \
  --no-cache \
  --security-opt no-new-privileges \
  -t payment:hardened \
  ~/payment/

sudo podman images payment

TASK 3 
sudo podman run -d \
  --name payment-instance \
  --network host \
  --memory 256m \
  --memory-swap 256m \
  --cpus 0.5 \
  --security-opt no-new-privileges \
  --cap-drop ALL \
  payment:hardened

sleep 5
curl -I http://localhost:4200/
sudo podman ps
sudo podman login docker.io
sudo podman tag payment:hardened docker.io/YOURUSERNAME/payment:hardened
sudo podman push docker.io/YOURUSERNAME/payment:hardened

TASK 4
sudo dnf install -y audit

sudo vi /usr/local/bin/build-pipeline.sh
#!/bin/bash
set -e
TIMESTAMP=$(date +%s)
CONTAINER_NAME="payment-test-${TIMESTAMP}"
WORK_DIR="/tmp/payment-build-${TIMESTAMP}"
IMAGE_NAME="payment:hardened"

# Step 1: git clone
git clone https://github.com/dvalyayevkbtu/payment "$WORK_DIR"

# Step 2: build
cd "$WORK_DIR"
sudo podman build \
  --no-cache \
  --security-opt no-new-privileges \
  -t "$IMAGE_NAME" .

# Step 3: start
sudo podman run -d \
  --name "$CONTAINER_NAME" \
  --network host \
  --memory 256m \
  --cpus 0.5 \
  --security-opt no-new-privileges \
  "$IMAGE_NAME"

# Step 4: check
sleep 5
curl -sf --max-time 10 http://localhost:4200 && echo "OK" || echo "FAILED"

# Step 5: delete container
sudo podman rm -f "$CONTAINER_NAME"

# Step 6: delete repo
rm -rf "$WORK_DIR"

echo "Pipeline complete"
bashsudo chown root:operations /usr/local/bin/build-pipeline.sh
sudo chmod 750 /usr/local/bin/build-pipeline.sh
sudo cp /usr/local/bin/build-pipeline.sh /home/kbtu/build-pipeline.sh

sudo vi /etc/audit/rules.d/pipeline.rules

-a always,exit -F path=/usr/local/bin/build-pipeline.sh -F perm=x -k pipeline-execution
sudo cp /etc/audit/rules.d/pipeline.rules /home/kbtu/pipeline.rules
sudo augenrules --load
sudo systemctl enable --now auditd

TASK 5 — nftables LAST
sudo vi /etc/nftables.conf

#!/usr/sbin/nft -f
flush ruleset
table inet filter {
  chain input {
    type filter hook input priority 0;
    ct state established,related accept
    iif "lo" accept
    tcp dport { 22, 4200 } accept
    drop
  }
  chain forward { type filter hook forward priority 0; }
  chain output { type filter hook output priority 0; }
}
sudo cp /etc/nftables.conf /home/kbtu/nftables.conf
sudo nft -f /etc/nftables.conf
sudo systemctl enable nftables
