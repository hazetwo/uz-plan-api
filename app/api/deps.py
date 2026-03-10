from typing import Annotated

import httpx
from fastapi import Depends, Request


def get_http_client(request: Request) -> httpx.AsyncClient:
    client: httpx.AsyncClient = request.app.state.http
    return client


HttpClient = Annotated[httpx.AsyncClient, Depends(get_http_client)]
