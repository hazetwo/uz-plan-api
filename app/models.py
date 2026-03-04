from datetime import date, time

from pydantic import BaseModel


class ScheduleEntry(BaseModel):
    date: date
    subGroup: str | None = None
    start: time | None = None
    end: time | None = None
    subject: str
    type: str
    teacher: list[str]
    room: str | None = None
