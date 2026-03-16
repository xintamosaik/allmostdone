import { Todo } from "./todo";
 
const todo = new Todo(1, {
    short: "Buy milk",
    description: "Semi-skimmed",
    due_date: "2026-03-20",
    cost_of_delay: 1,
    effort: "hours",
});

console.log(todo.renderTableRow());
console.log(todo.renderCard());
console.log(todo.renderEditForm("/todos/1/update"));
console.log(todo.toJson());
 

todo.apply({
    short: "Buy bread",
    description: "Sourdough",
    due_date: "2026-03-21",
    cost_of_delay: "2",
    effort: "days",
});

console.log(todo.renderTableRow());
console.log(todo.toJson());