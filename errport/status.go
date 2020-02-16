package errport

type status string

const StatusFixed status = "fixed"
const StatusForReview status = "for_review"
const StatusIgnored status = "ignored"
const StatusInProgress status = "in progress"
const StatusOpen status = "open"
const StatusSnoozed status = "snoozed"
