# Prompts Used — AI Collaboration Log

This document logs the prompts used with Claude throughout the
development of this project. Prompts are curated (trivial
confirmations like "ok" or "looks good" are omitted) and lightly
copy-edited for grammar and structure only — the actual requests,
questions, and technical decisions are unchanged from the original
conversation. The full, unedited conversation is available on
request.

The goal of this log is to be an honest artifact of how the project
was actually built: iterative, with real back-and-forth, corrections,
and course changes — not a polished transcript written after the
fact.

## Phase 1 — Project kickoff and learning approach

> Help me build this project [full challenge brief pasted: full-stack
> calculator, React frontend, Go backend microservice, REST API,
> clean architecture, unit tests, documentation]. I want us to do
> this step by step so I understand how everything works. To start:
> I have never worked with Go before.

> Wait, you're moving too fast. I don't even know which IDE to use to
> see my Go code, where the project folder will be saved, and we
> haven't created the GitHub repository yet — none of that is set up.

> I want us to move step by step, efficiently but calmly, where you
> and I decide together how to do things — you tell me when to write
> code, and you explain what was done. I also want us to choose the
> architecture, design patterns, best practices, entities, and
> relationships very carefully. And I want everything documented
> professionally, since this is how I'll be presenting myself for a
> job and I need to get accepted.

## Phase 2 — Architecture decisions

> We could add persistence to make it more impressive. I also love
> Clean Architecture, hexagonal architecture, SOLID, and design
> patterns in general.

> [Pasted the full Sezzle Software Engineer Intern job posting,
> including required tech stack: Golang, TypeScript/React, Postgres,
> Git, GitLab CI/CD, and the explicit requirement of demonstrated
> experience working with Claude or an equivalent LLM.]

> Which database do you think is best, and why — considering ease of
> deployment and handoff to Sezzle's team?

> Is it fine for us to push everything to `main`? Or should we have
> `main`, `develop`, and a branch per feature?

> Shouldn't the full hexagonal architecture have a few more pieces?
> (referring to a `usecase`-style layer)

## Phase 3 — Debugging and working method

> I don't understand what I need to change or why — give me the
> complete, corrected file.
> *(Context: preferred receiving the full corrected file over a
> partial diff, to avoid ambiguity when editing code directly.)*

> [Pasted a terminal error after a failed `docker compose up`,
> including the full build log] — from where do I run this?

## Phase 4 — Software architecture research

> Can you search the internet for the best ways to design software
> architecture, so we can make this project better and more
> professional?

> But we don't have a `usecase` layer — does that mean something is
> missing?

## Phase 5 — Frontend design

> I want you to research the most current and trending frontend
> styles, and design something that doesn't look like the generic
> style AI usually produces — like typical corporate landing pages.

> I've seen a lot of "liquid glass" and minimalism trends lately —
> let's combine those with the trends you already found into
> something visually strong.

## Phase 6 — Reviewing and correcting AI output

> [After reviewing a deployment diagram Claude generated] Are you
> sure this is wrong? What I'm trying to represent is that the user
> doesn't know anything about how the app works internally — they
> just connect through their laptop, go to the internet, reach the
> frontend, and the frontend is what calls the backend to consume its
> services. Isn't that right?

> I understand, but how could this actually be shown in the diagram —
> that the user goes to the frontend, the frontend only serves the
> SPA, and everything else is handled by calls made from the browser
> itself?

> [After reviewing a self-made package diagram] Here's another one I
> built — tell me honestly if it's correct (communications,
> protocols, direction of dependencies, etc.).

> I don't want to include generic diagrams made by you, like this
> one — make it yourself, and make sure it's actually clear this
> time, because the way it is now, I don't understand it.

## Phase 7 — Documentation quality

> Since what I'm trying to prove is that I'm skilled and that I
> handle AI tools well — since that matters a lot for this role —
> I want the prompts we share to be genuinely good too. I know some
> of mine weren't great. So include the list, but skip the ones that
> don't add value or are too redundant, and clean up the writing on
> the ones we do include so they read as clear, structured, and
> professional — without fabricating skills or decisions that weren't
> actually mine.

## Phase 8 — Diagram rendering troubleshooting

> [After noticing the architecture diagrams weren't displaying in the
> GitHub README] I want to review what happened here.

> [After a suggested fix only partially worked] I don't care how —
> I need us to fix this so all of them render correctly under any
> theme, the same way the others already do.

> Let's just make it a PNG and be done with it.
> *(Context: after several rounds of trying to fix an SVG rendering
> issue at the root cause, chose the simpler, more robust option
> instead of continuing to chase an edge case.)*

> Can we push this to the same branch?
> *(Context: asked before reusing an already-merged branch for an
> unrelated fix — checking Git hygiene before acting, not just
> following instructions blindly.)*

## Phase 9 — API documentation (Swagger)

> Is Swagger compatible with Go for documenting APIs?

> Yes, let's add it.

---

**Note on methodology:** every architectural decision in this
project (hexagonal architecture, the Strategy pattern for operations,
using two repository adapters, GitHub Flow over Git Flow, CSS Modules
over a UI library, no state-management library on the frontend) was
discussed and reasoned through in conversation before being
implemented — none were accepted as a default suggestion without
understanding the trade-off first.
