import pytest

from app.core.scraper import fetch_groups, fetch_schedule


@pytest.mark.asyncio
async def test_fetch_schedule(client):
    response = await fetch_schedule("https://www.python.org/", client)
    assert response is not None


@pytest.mark.asyncio
async def test_fetch_groups(client):
    response = await fetch_groups(client)
    assert response is not None
