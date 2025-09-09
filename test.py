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
    print("Registration successful", resp.status_code)

    resp = session.post(URL + "/login", json={'login': USER_LOGIN, 'password': USER_PASSWORD})
    if resp.status_code != 200:
        print('Login failed:', resp.text)
        return
    print("Login successful", resp.status_code, resp.headers, resp.text)

    resp = session.get(URL + "/healthz")
    print("health", resp.request.headers)


if __name__ == "__main__":
    test_api()