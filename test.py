import requests

URL = "http://localhost:8080"

USER_LOGIN = 'test_user'
USER_PASSWORD = 'test_password'

def test_api():
    session = requests.Session()
    resp = session.post(URL + "/check_token", json={'login': USER_LOGIN, 'jwt': "some_jwt_token"})
    if resp.status_code != 401:
        print("Check token failed: expected 401, got", resp.status_code)
        return
    else:
        print("OK: Check token unauthorized")

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
    coockies = session.cookies.get_dict()
    if resp.status_code != 200:
        print("Health check failed:", resp, resp.text, resp.headers)
        return
    elif coockies.get("jwt") is None:
        print("Health check failed: no jwt coockie")
        print(resp.request.headers)
        return
    else:
        print("OK: health")
    resp = session.get(URL + "/check_token")
    if resp.status_code != 200:
        print("Check token failed:", resp.status_code, resp.text)
        return
    else:
        print("OK: Check token authorized")


if __name__ == "__main__":
    test_api()