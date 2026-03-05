import httpx


def create_http_client():
    timeout = httpx.Timeout(timeout=10)
    return httpx.AsyncClient(timeout=timeout, http2=True)
