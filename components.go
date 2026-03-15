package main

const backButton = `
{{define "backButton"}}
<button
  type="button"
  fx-action="/todos/list"
  fx-target="#output"
  fx-swap="innerHTML">
  Back
</button>
{{end}}
`