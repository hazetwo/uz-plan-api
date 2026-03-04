import httpx
from fastapi import Request


def create_http_client():
    timeout = httpx.Timeout(timeout=10)
    return httpx.AsyncClient(timeout=timeout, http2=True)


def get_http_client(request: Request) -> httpx.AsyncClient:
    return request.app.state.http
