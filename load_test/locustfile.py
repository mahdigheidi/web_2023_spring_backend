from locust import HttpUser, task

class BackendLoadTest(HttpUser):
    @task
    def load_test(self):
        self.client.get("/auth/req_pq", verify=False)
        self.client.get("/auth/req_dh_params", verify=False)
        self.client.get("/biz/get_users", verify=False)
        self.client.get("/biz/get_users_with_sql_injection", verify=False)