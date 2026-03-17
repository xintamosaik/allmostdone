import { Todo, type TodoRawInput } from "./todo";

/**
 * A simple list with todo items and an action button to edit them
 */
function TodoList() {
    const rows = todos.map(todo => todo.renderTableRow());
    return `
        <button
            type="button"
            fx-action="/todos/new"
            fx-target="#output"
            fx-swap="innerHTML">
            New Todo
        </button>
        <table>
            <thead>
                <tr>
                    <th>ID</th>
                    <th>Short</th>
                    <th>Description</th>
                    <th>Due Date</th>
                    <th>Cost of Delay</th>
                    <th>Effort</th>
                    <th>Actions</th>
                </tr>
            </thead>
            <tbody>
                ${rows.join("\n")}
            </tbody>
        </table>
    `;
}

/**
 * A form to edit a todo item
 */
function EditTodo(id: string) {
    if (!id) {
        return htmlResponse("Not Found", 404);
    }

    const todo = todos.find(t => t.id() === Number(id));
    if (!todo) {
        return htmlResponse("Not Found", 404);
    }

    return htmlResponse(`Edit todo ${id}: ${todo.renderEditForm()}`);
}

function CreateTodo() {
    const newTodo = new Todo(todos.length + 1, {
        short: "...",
        description: "",
        effort: 'hours',
        cost_of_delay: 0,
        due_date: new Date().toISOString().split("T")[0],
    });
    todos.push(newTodo);
    return htmlResponse(`Create new todo: ${newTodo.renderEditForm()}`);
}

/**
 * Parses the form data and updates the todo item. If there are any validation errors, it returns a 400 response with the error messages. Otherwise, it returns the updated todo list.
 */
async function parseEdit( req: Bun.BunRequest) {
    const id = req.params["id"];
    if (!id) {
        return htmlResponse("Not Found", 404);
    }
    const todo = todos.find(t => t.id() === Number(id));
    if (!todo) {
        return htmlResponse("Not Found", 404);
    }

    const formData = await req.formData();
    const data = {
        short: formData.get("short") as string,
        description: formData.get("description") as string,
        effort: formData.get("effort") as string,
        cost_of_delay: Number(formData.get("cost_of_delay")),
        due_date: formData.get("due_date") as string,
    } as TodoRawInput;

    const result = todo.apply(data);
    if (!result.ok) {
        return htmlResponse(`Error: ${result.errors.join(", ")}`, 400);
    }

    return htmlResponse(TodoList(), 200);
}
const todos: Todo[] = [];
const example: TodoRawInput = {
    short: "Example Todo",
    description: "This is an example todo",
    effort: 'hours',
    cost_of_delay: 1,
    due_date: "2024-12-31",
}

function htmlResponse(html: string, status = 200): Response {
    return new Response(html, {
        status,
        headers: {
            "Content-Type": "text/html; charset=utf-8",
        },
    });
}

const exampleTodo = new Todo(1, example)
todos.push(exampleTodo);

type AppRoute =
    | "/"
    | "/fixi-0.9.2.js"
    | "/style.css"
    | "/status"
    | "/todos/list"
    | "/todos/:id/edit"
    | "/todos/:id/update"
    | "/todos/new"
    | "/favicon.ico";

const routes = {
    // INDEX
    "/": Bun.file("./index.html"),

    // FIXI
    "/fixi-0.9.2.js": Bun.file("./static/fixi-0.9.2.js"),

    // CSS
    "/style.css": Bun.file("./static/style.css"),

    // STATUS
    "/status": htmlResponse("OK"),

    // LIST
    "/todos/list": () => htmlResponse(TodoList()),

    // EDIT
    "/todos/:id/edit": (req: Bun.BunRequest<"/todos/:id/edit">) => EditTodo(req.params["id"]),

    // UPDATE
    "/todos/:id/update": {
        POST: async (req: Bun.BunRequest<"/todos/:id/update">) => parseEdit(req),
    },

    // CREATE
    "/todos/new": () => CreateTodo(),

    // FAVICON
    "/favicon.ico": Bun.file("./favicon.ico"),
} satisfies Bun.Serve.Routes<undefined, AppRoute>;

const server = Bun.serve({
    routes,

    // CATCH ALL
    fetch(_: Request) {
        return new Response("Not Found", { status: 404 });
    },
});

console.log(`Server running at ${server.url}`);
