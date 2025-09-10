import requests

URL = "http://localhost:8080"

USER_LOGIN = 'test_user'
USER_PASSWORD = 'test_password'

def test_api():
    session = requests.Session()
    resp = session.post(URL + "/reg", json={'login': USER_LOGIN, 'password': USER_PASSWORD})
    if resp.status_code != 200:
        print("Registration failed:", resp.text)
        return
    else:
        print("OK: Registration")

    resp = session.post(URL + "/login", json={'login': USER_LOGIN, 'password': USER_PASSWORD})
    if resp.status_code != 200:
        print('Login failed:', resp.text)
        return
    else:
        print("OK: Login successful")

    resp = session.get(URL + "/healthz")
    if resp.status_code != 200:
        print("Health check failed:", resp, resp.text, resp.headers)
        return
    else:
        print("OK: health")


if __name__ == "__main__":
    test_api()