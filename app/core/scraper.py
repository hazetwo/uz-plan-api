import httpx
from bs4 import BeautifulSoup

from app.core.exceptions import FetchScheduleException, UrlNotFoundException


async def async_fetch(
    url: httpx.URL, client: httpx.AsyncClient
) -> BeautifulSoup:
    try:
        response = await client.get(url)
        response.raise_for_status()
    except httpx.RequestError as exc:
        raise FetchScheduleException(f"Request error: {exc}")
    except httpx.HTTPStatusError as exc:
        raise FetchScheduleException(f"HTTP error: {exc}")

    return BeautifulSoup(response.text, "html.parser")


async def fetch_schedule(url: str, client: httpx.AsyncClient) -> BeautifulSoup:
    if url is None:
        raise UrlNotFoundException(
            "The provided ID does not map to a valid URL."
        )

    soup = await async_fetch(httpx.URL(url), client)

    return soup
