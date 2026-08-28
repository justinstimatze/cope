package scan

// The lane for prose that leaves the conversation.
//
// Everything else this package scores is a turn: something written to a person
// who is about to answer it. A ticket description, a project brief, a status
// update is written to somebody who will read it days later with none of the
// session around it, and who cannot answer "continue" because there is no turn
// to continue. Four rules assume that reader and misfire here — see
// needsAnswerableReader — and a scorer that flags dangling_end on every ticket
// is a scorer nobody leaves switched on.
//
// Nothing is added in exchange. The loop lane swapped its four for two, because
// an unattended report has a failure of its own to check. External prose does
// not: what is left after the drop — the voicing rules, labelled_opening,
// clause_symmetry, paragraph_uniformity — is the whole of what this lane is
// for, and it is more load on those rules rather than less. A ticket is read
// cold by definition, which is the condition every one of them was written for.
const LaneExternal = "external"
