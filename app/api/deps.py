from typing import Annotated

import httpx
from fastapi import Depends, Request

from app.models import Group


def get_http_client(request: Request) -> httpx.AsyncClient:
    client: httpx.AsyncClient = request.app.state.http
    return client


def get_groups_data(request: Request) -> list[Group]:
    groups_data: list[Group] = request.app.state.groups
    return groups_data


HttpClient = Annotated[httpx.AsyncClient, Depends(get_http_client)]
Groups = Annotated[list[Group], Depends(get_groups_data)]
