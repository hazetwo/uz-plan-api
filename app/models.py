from datetime import date, time

from pydantic import BaseModel


class ScheduleEntry(BaseModel):
    time: date
    subGroup: str | None = None
    start: time
    end: time
    subject: str
    type: str
    teacher: list[str]
    room: str | None = None
