from datetime import date
from typing import List

from fastapi import APIRouter

from app.api.deps import Groups, HttpClient
from app.core.parser import parse_schedule
from app.core.scraper import fetch_schedule
from app.models import ScheduleEntry
from app.utils.date import get_week_end

router = APIRouter(prefix="/schedule", tags=["schedule"])


async def get_schedule(
    id: str, groups_data: Groups, client: HttpClient
) -> List[ScheduleEntry]:
    soup = await fetch_schedule(id, groups_data, client)
    schedules_entries = parse_schedule(soup)
    # Push date=null to the bottom of the list
    schedules_entries.sort(key=lambda entry: (entry.date is None, entry.date))
    return schedules_entries


@router.get("/")
async def get_entries(
    id: str, groups_data: Groups, client: HttpClient
) -> List[ScheduleEntry]:
    return await get_schedule(id, groups_data, client)


@router.get("/by-day")
async def get_entry_by_day(
    id: str, groups_data: Groups, date: date, client: HttpClient
) -> List[ScheduleEntry]:
    schedule_entries = await get_schedule(id, groups_data, client)
    day_entries: List[ScheduleEntry] = [
        entry for entry in schedule_entries if date == entry.date
    ]

    return day_entries


@router.get("/by-week")
async def get_entries_by_week(
    id: str, date: date, groups_data: Groups, client: HttpClient
) -> List[ScheduleEntry]:
    schedule_entries = await get_schedule(id, groups_data, client)
    week_end = get_week_end(date)
    week_entries: List[ScheduleEntry] = [
        entry
        for entry in schedule_entries
        if entry.date is not None and date <= entry.date <= week_end
    ]

    return week_entries
