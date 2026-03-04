import pytest

from app.api.core.scraper import fetch_schedule


@pytest.mark.asyncio
async def test_fetch_schedule(client):
    response = await fetch_schedule("https://www.python.org/", client)
    assert response is not None
