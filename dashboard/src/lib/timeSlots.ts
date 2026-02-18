import Moment from "moment";
import { extendMoment } from "moment-range";

const moment = extendMoment(Moment);

export function timestampRange(
  data: Data,
  token: string | null
): [Date, Date] | null {
  if (token === null) return null;

  let max = new Date(-8640000000000000);
  let min = new Date(8640000000000000);
  let found = false;

  for (const extraction of data.extraction) {
    if (!(token in extraction.params)) continue;
    const ts = new Date(extraction.params[token].value as string);
    if (ts > max) max = ts;
    if (ts < min) min = ts;
    found = true;
  }

  return found ? [min, max] : null;
}

export function makeTimeSlots(
  min: Date,
  max: Date,
  stepMinutes: number
): Moment.Moment[] {
  return Array.from(
    moment.range(min, max).by("minutes", { step: stepMinutes })
  );
}

export function nearestSlotIndex(
  timestampMs: number,
  slotMs: number[]
): number {
  let bestIndex = 0;
  let bestDiff = Math.abs(slotMs[0] - timestampMs);
  for (let j = 1; j < slotMs.length; j++) {
    const diff = Math.abs(slotMs[j] - timestampMs);
    if (diff < bestDiff) {
      bestIndex = j;
      bestDiff = diff;
    } else {
      // Slots are sorted; once diff grows we've passed the closest slot
      break;
    }
  }
  return bestIndex;
}
