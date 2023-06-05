from locust import HttpUser, task

class BackendLoadTest(HttpUser):
    @task
    def load_test(self):
        self.client.get("/auth/req_pq?nonce=ABCDE01234ABCDE01237&message_id=0", verify=False)
        self.client.get("/auth/req_dh_params?nonce=ABCDE01234ABCDE01237&server_nonce=93FdrQPTRkl897Etrw7T&message_id=4&a=8", verify=False)
        self.client.get("/biz/get_users?user_id=2&message_id=2&auth_key=2&auth_id=ABCDE01234ABCDE0123793FdrQPTRkl897Etrw7T", verify=False)
        self.client.get("/biz/get_users?user_id=mahdigheidi&message_id=2&auth_key=2&auth_id=ABCDE01234ABCDE0123793FdrQPTRkl897Etrw7T", verify=False)