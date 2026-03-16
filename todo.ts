/**
 * The idea of this file is to completely represent a Todo item. 
 * 
 * We deliberately include persistence and representation to make this a deep module aka Ousterhouts design.
 * 
 * This is just an experiment to compare different styles of development. 
 * 
 * Anyone working on this is not allowed to change the philosophy of this module but is free to change the implementation.
 * 
 * What is encouraged: Add MongoDB in the same style as SQL was added. Add JSON representation for non-HTML consuming Frontends
 * 
 * Add other forms of persistence or representation as you see fit. 
 * 
 * Also encouraged: Try to be LESS generic. Re-usability is a non-concern here. We go low abstraction WITHIN the file WHILE abstracting for the consumer.
 * 
 * Also really encouraged: A lot of comments on methods and classes. 
 * 
 * NOT encouraged: Adding a "TodoList" class or similar. This file is about the Todo item, not about collections of Todos.
 */

type SqlParam = string | number | null;

type TodoInitial = {
    short?: string;
    description?: string;
    due_date?: string;
    cost_of_delay?: number;
    effort?: string;
};

type TodoRawInput = {
    short?: string;
    description?: string;
    due_date?: string;
    cost_of_delay?: string;
    effort?: string;
};

/**
 * JSON payload kept close to the Todo model so representation rules stay in one place.
 */
type TodoJson = {
    id: number;
    short: string;
    description: string;
    due_date: string | null;
    cost_of_delay: number;
    effort: string;
};

/**
 * Dedicated field for the todo title. Required and intentionally strict.
 */
class TodoShort {
    private _value: string;

    constructor(initialValue: string) {
        this._value = "";
        this.setFromRaw(initialValue);
    }

    setFromRaw(raw: string): void {
        const cleaned = (raw ?? "").trim();

        if (cleaned.length === 0) {
            throw new Error("Short must not be blank");
        }

        if (cleaned.length > 120) {
            throw new Error("Short must be <= 120 chars");
        }

        this._value = cleaned;
    }

    valueAsString(): string {
        return this._value;
    }

    valueAsSqlParam(): string {
        return this._value;
    }

    valueAsJson(): string {
        return this._value;
    }

    renderField(): string {
        return `
      <div>
        <label for="short">Short</label>
        <input
          type="text"
          id="short"
          name="short"
          value="${escapeHtml(this._value)}"
        >
      </div>
    `.trim();
    }

    renderTableCell(): string {
        return `<td>${escapeHtml(this._value)}</td>`;
    }
}

/**
 * Dedicated field for detailed description. Optional but bounded.
 */
class TodoDescription {
    private _value: string;

    constructor(initialValue: string) {
        this._value = "";
        this.setFromRaw(initialValue);
    }

    setFromRaw(raw: string): void {
        const cleaned = (raw ?? "").trim();

        if (cleaned.length > 5000) {
            throw new Error("Description must be <= 5000 chars");
        }

        this._value = cleaned;
    }

    valueAsString(): string {
        return this._value;
    }

    valueAsSqlParam(): string {
        return this._value;
    }

    valueAsJson(): string {
        return this._value;
    }

    renderField(): string {
        return `
      <div>
        <label for="description">Description</label>
        <textarea id="description" name="description">${escapeHtml(this._value)}</textarea>
      </div>
    `.trim();
    }

    renderTableCell(): string {
        return `<td>${escapeHtml(this._value)}</td>`;
    }
}

/**
 * Dedicated due date with a narrow accepted format for consistency with SQL and forms.
 */
class TodoDueDate {
    private _value: string;

    constructor(initialValue: string) {
        this._value = "";
        this.setFromRaw(initialValue);
    }

    setFromRaw(raw: string): void {
        const cleaned = (raw ?? "").trim();

        if (cleaned === "") {
            this._value = "";
            return;
        }

        if (!/^\d{4}-\d{2}-\d{2}$/.test(cleaned)) {
            throw new Error("Due Date must be YYYY-MM-DD");
        }

        this._value = cleaned;
    }

    valueAsString(): string {
        return this._value;
    }

    valueAsSqlParam(): string | null {
        return this._value === "" ? null : this._value;
    }

    valueAsJson(): string | null {
        return this._value === "" ? null : this._value;
    }

    renderField(): string {
        return `
      <div>
        <label for="due_date">Due Date</label>
        <input
          type="date"
          id="due_date"
          name="due_date"
          value="${escapeHtml(this._value)}"
        >
      </div>
    `.trim();
    }

    renderTableCell(): string {
        return `<td>${escapeHtml(this._value)}</td>`;
    }
}

/**
 * Domain-specific integer: cost of delay for this todo, constrained to a small scale.
 */
class TodoCostOfDelay {
    private _value: number;

    constructor(initialValue: number) {
        this._value = 0;
        this.setFromNumber(initialValue);
    }

    setFromRaw(raw: string): void {
        const cleaned = (raw ?? "").trim();
        const parsed = Number.parseInt(cleaned, 10);

        if (Number.isNaN(parsed)) {
            throw new Error("Cost Of Delay must be an integer");
        }

        this.setFromNumber(parsed);
    }

    valueAsString(): string {
        return String(this._value);
    }

    valueAsSqlParam(): number {
        return this._value;
    }

    valueAsJson(): number {
        return this._value;
    }

    renderField(): string {
        return `
      <div>
        <label for="cost_of_delay">Cost Of Delay</label>
        <input
          type="number"
          id="cost_of_delay"
          name="cost_of_delay"
          value="${this._value}"
          min="-2"
          max="2"
        >
      </div>
    `.trim();
    }

    renderTableCell(): string {
        return `<td>${this._value}</td>`;
    }

    private setFromNumber(value: number): void {
        if (value < -2) {
            throw new Error("Cost Of Delay must be >= -2");
        }

        if (value > 2) {
            throw new Error("Cost Of Delay must be <= 2");
        }

        this._value = value;
    }
}

/**
 * Domain-specific selection for effort sizing.
 */
class TodoEffort {
    private _value: string;
    private _options: string[];

    constructor(initialValue: string) {
        this._value = "";
        this._options = ["mins", "hours", "days", "weeks", "months"];
        this.setFromRaw(initialValue);
    }

    setFromRaw(raw: string): void {
        const cleaned = (raw ?? "").trim();

        if (!this._options.includes(cleaned)) {
            throw new Error(`Invalid effort: ${cleaned}`);
        }

        this._value = cleaned;
    }

    valueAsString(): string {
        return this._value;
    }

    valueAsSqlParam(): string {
        return this._value;
    }

    valueAsJson(): string {
        return this._value;
    }

    renderField(): string {
        const optionsHtml = this._options
            .map((option) => {
                const selected = option === this._value ? " selected" : "";
                return `<option value="${escapeHtml(option)}"${selected}>${escapeHtml(option)}</option>`;
            })
            .join("");

        return `
      <div>
        <label for="effort">Effort</label>
        <select id="effort" name="effort">
          ${optionsHtml}
        </select>
      </div>
    `.trim();
    }

    renderTableCell(): string {
        return `<td>${escapeHtml(this._value)}</td>`;
    }
}
/**
 * The Todo class. Some would call it a god-object, but I view it as a deep module (for now). I might changem my view.
 * 
 * The Todo class should enable users to create Todo items and change them safely. It should enable them to persist them without issues and to render representations and forms in an easy manner.
 */
class Todo {
    private _id: number;

    private shortField: TodoShort;
    private descriptionField: TodoDescription;
    private dueDateField: TodoDueDate;
    private costOfDelayField: TodoCostOfDelay;
    private effortField: TodoEffort;

    constructor(id: number, initial?: TodoInitial) {
        this._id = id;

        this.shortField = new TodoShort(initial?.short ?? "");
        this.descriptionField = new TodoDescription(initial?.description ?? "");
        this.dueDateField = new TodoDueDate(initial?.due_date ?? "");
        this.costOfDelayField = new TodoCostOfDelay(initial?.cost_of_delay ?? 0);
        this.effortField = new TodoEffort(initial?.effort ?? "hours");
    }

    id(): number {
        return this._id;
    }

    short(): string {
        return this.shortField.valueAsString();
    }

    description(): string {
        return this.descriptionField.valueAsString();
    }

    dueDate(): string {
        return this.dueDateField.valueAsString();
    }

    costOfDelay(): number {
        return Number(this.costOfDelayField.valueAsSqlParam());
    }

    effort(): string {
        return this.effortField.valueAsString();
    }

    apply(raw: TodoRawInput): void {
        // Apply all user-provided values as one transaction-like operation.
        this.shortField.setFromRaw(raw.short ?? "");
        this.descriptionField.setFromRaw(raw.description ?? "");
        this.dueDateField.setFromRaw(raw.due_date ?? "");
        this.costOfDelayField.setFromRaw(raw.cost_of_delay ?? "");
        this.effortField.setFromRaw(raw.effort ?? "");
    }

    renderTableRow(): string {
        return `
      <tr>
        <td>${this._id}</td>
        ${this.shortField.renderTableCell()}
        ${this.descriptionField.renderTableCell()}
        ${this.dueDateField.renderTableCell()}
        ${this.costOfDelayField.renderTableCell()}
        ${this.effortField.renderTableCell()}
      </tr>
    `.trim();
    }

    renderCard(): string {
        return `
      <article class="todo-card">
        <h3>${escapeHtml(this.short())}</h3>
        <p>${escapeHtml(this.description())}</p>
        <dl>
          <div><dt>Due date</dt><dd>${escapeHtml(this.dueDate() || "N/A")}</dd></div>
          <div><dt>Cost of delay</dt><dd>${this.costOfDelay()}</dd></div>
          <div><dt>Effort</dt><dd>${escapeHtml(this.effort())}</dd></div>
        </dl>
      </article>
    `.trim();
    }

    renderEditForm(action: string): string {
        const fieldsHtml = [
            this.shortField.renderField(),
            this.descriptionField.renderField(),
            this.dueDateField.renderField(),
            this.costOfDelayField.renderField(),
            this.effortField.renderField(),
        ].join("\n");

        return `
      <form method="post" action="${escapeHtml(action)}">
        ${fieldsHtml}
        <div>
          <input type="submit" value="Save">
        </div>
      </form>
    `.trim();
    }

    toJson(): TodoJson {
        return {
            id: this._id,
            short: this.shortField.valueAsJson(),
            description: this.descriptionField.valueAsJson(),
            due_date: this.dueDateField.valueAsJson(),
            cost_of_delay: this.costOfDelayField.valueAsJson(),
            effort: this.effortField.valueAsJson(),
        };
    }

    insertSql(tableName = "todos"): { sql: string; params: SqlParam[] } {
        const params: SqlParam[] = [
            this.shortField.valueAsSqlParam(),
            this.descriptionField.valueAsSqlParam(),
            this.dueDateField.valueAsSqlParam(),
            this.costOfDelayField.valueAsSqlParam(),
            this.effortField.valueAsSqlParam(),
        ];

        return {
            sql: `INSERT INTO ${tableName} (short, description, due_date, cost_of_delay, effort) VALUES ($1, $2, $3, $4, $5)`,
            params,
        };
    }

    updateSql(tableName = "todos"): { sql: string; params: SqlParam[] } {
        const params: SqlParam[] = [
            this.shortField.valueAsSqlParam(),
            this.descriptionField.valueAsSqlParam(),
            this.dueDateField.valueAsSqlParam(),
            this.costOfDelayField.valueAsSqlParam(),
            this.effortField.valueAsSqlParam(),
        ];

        return {
            sql: `UPDATE ${tableName} SET short = $1, description = $2, due_date = $3, cost_of_delay = $4, effort = $5 WHERE id = $6`,
            params: [...params, this._id],
        };
    }
}

function escapeHtml(value: string): string {
    return value
        .replaceAll("&", "&amp;")
        .replaceAll("<", "&lt;")
        .replaceAll(">", "&gt;")
        .replaceAll('"', "&quot;")
        .replaceAll("'", "&#39;");
}

export { Todo };