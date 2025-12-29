# JSON Parser

This is a JSON parser written from scratch in Go . It aims to be fully compatible 
with [RFC8259](https://datatracker.ietf.org/doc/html/rfc8259) . 
<details>
  <summary>Automata to identify numbers</summary>
  <img width="564" height="336" alt="automata" src="https://github.com/user-attachments/assets/fd695114-f66e-4814-8ede-1016c24148dc" />
  <br>To identify RFC8259 compliant numbers. It simulates the automata given above. Every character that doesn't have a transition leads to 
  a dead state . Since the language contains every UTF-8 character , it is not possible to show everything . 
  <br>
  <h3>AI Disclaimer</h3>
  I used NanoBanana to convert the hand‑drawn automaton into a digital illustration.
</details>
