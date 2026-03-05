from datetime import date
from typing import List

from fastapi import APIRouter

from app.api.deps import HttpClient
from app.core.parser import parse_schedule
from app.core.scraper import fetch_schedule
from app.models import ScheduleEntry
from app.utils import get_week_end

router = APIRouter(prefix="/schedule", tags=["items"])


@router.get("/")
async def get_entries(id: str, client: HttpClient) -> List[ScheduleEntry]:
    soup = await fetch_schedule(id, client)
    schedules_entries = parse_schedule(soup)

    return schedules_entries


@router.get("by-day")
async def get_entry(id: str, date: date, client: HttpClient):
    soup = await fetch_schedule(id, client)
    schedule_entries = parse_schedule(soup)
    day_entries: List[ScheduleEntry] = []
    for entry in schedule_entries:
        if date == entry.date:
            day_entries.append(entry)

    return day_entries


@router.get("/by-week")
async def get_week_entries(
    id: str, date: date, client: HttpClient
) -> List[ScheduleEntry]:
    soup = await fetch_schedule(id, client)
    schedule_entries = parse_schedule(soup)
    week_entries: List[ScheduleEntry] = []
    week_end = get_week_end(date)

    for entry in schedule_entries:
        if date <= entry.date <= week_end:
            week_entries.append(entry)

    return week_entries
