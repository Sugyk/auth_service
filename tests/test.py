import requests
import jwt

URL = "http://localhost:8080"

USER_LOGIN = 'test_user'
USER_PASSWORD = 'test_password'

def test_api():
    session = requests.Session()
    resp = session.post(URL + "/check_token", json={'login': USER_LOGIN})
    if resp.status_code != 400:
        print("Check token failed: expected 400, got", resp.status_code)
        return
    else:
        print("OK: Check token unauthorized")

    resp = session.post(URL + "/reg", json={'login': USER_LOGIN, 'password': USER_PASSWORD})
    if resp.status_code != 200:
        print("Registration failed:", resp.text)
        return
    else:
        print("OK: Registration")

    resp = session.post(URL + "/reg", json={'login': USER_LOGIN, 'password': USER_PASSWORD})
    if resp.status_code == 409:
        print("OK: Re-registration restricted")
    else:
        print("Re-registration is successed - fail:", resp.status_code, resp.text)
        return

    resp = session.post(URL + "/login", json={'login': USER_LOGIN, 'password': USER_PASSWORD})
    if resp.status_code != 200:
        print('Login failed:', resp.text)
        return
    else:
        print("OK: Login successful")

    resp = session.get(URL + "/healthz")
    jwt_token = session.cookies.get('jwt')

    try:
        jwt.decode(jwt_token, "", algorithms=['HS256'])
        isEmptySecretKey = True
    except jwt.InvalidTokenError:
        isEmptySecretKey = False

    if resp.status_code != 200:
        print("Health check failed:", resp, resp.text, resp.headers)
        return
    elif jwt_token is None:
        print("Health check failed: no jwt coockie")
        print(resp.request.headers)
        return
    elif isEmptySecretKey:
        print("Health check failed: jwt secretkey is empty")
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