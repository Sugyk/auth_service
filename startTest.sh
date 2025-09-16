sudo docker compose down > /dev/null
sudo docker compose build > /dev/null
sudo docker compose up -d > /dev/null
python3 tests/test.py