import { Todo, type TodoInitial } from "./todo";
 
// We use fixi on the frontend to do "hypermedia" and "HATEOS"

// ENDPOINTS
// GET /todos -> list of todos
// GET /todos/1 -> one todo
// POST /todos -> create a new todo
// PATCH /todos/1 -> update a todo
// DELETE /todos/1 -> delete a todo

/** 
 * On the frontend we only use that to initialize the web app: 
 * 
 * <script>
    const output = document.getElementById('output');
    fetch("/todos/list")
        .then(response => response.text())
        .then(html => output.innerHTML = html)
        .catch(error => console.error('Error fetching initial content:', error));
    </script>
 */

// we have to make sure to react to /todos/list first

const todos = [];
const example: TodoInitial = {
    short: "Example Todo",
    description: "This is an example todo",
    effort: 'hours',
    cost_of_delay: 1,
    due_date: "2024-12-31",
}
const exampleTodo = new Todo(1, example)
console.log(exampleTodo.values())
todos.push(exampleTodo);
const server = Bun.serve({
  // `routes` requires Bun v1.2.3+
  routes: {
    // Static routes
    "/api/status": new Response("OK"),

    // Dynamic routes
    "/users/:id": req => {
      return new Response(`Hello User ${req.params.id}!`);
    },

    // Per-HTTP method handlers
    "/api/posts": {
      GET: () => new Response("List posts"),
      POST: async req => {
        const body = await req.json();
        return Response.json({ created: true, ...body });
      },
    },

    // Wildcard route for all routes that start with "/api/" and aren't otherwise matched
    "/api/*": Response.json({ message: "Not found" }, { status: 404 }),

    // Redirect from /blog/hello to /blog/hello/world
    "/blog/hello": Response.redirect("/blog/hello/world"),

    // Serve a file by lazily loading it into memory
    "/favicon.ico": Bun.file("./favicon.ico"),
  },

  // (optional) fallback for unmatched routes:
  // Required if Bun's version < 1.2.3
  fetch(req) {
    return new Response("Not Found", { status: 404 });
  },
});

console.log(`Server running at ${server.url}`);