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
    short: string;
    description: string;
    due_date: string;
    cost_of_delay: string;
    effort: string;
};

type TodoPatchInput = {
    short?: string;
    description?: string;
    due_date?: string;
    cost_of_delay?: string;
    effort?: string;
};

type TodoValidationError = {
    field: string;
    message: string;
};

type TodoValidationResult = {
    ok: boolean;
    errors: TodoValidationError[];
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
    readonly _name = 'short';
    private _value: string;

    constructor(initialValue: string) {
        this._value = "";
        const error = this.setFromRaw(initialValue);
        if (error) {
            throw error;
        }
    }

    setFromRaw(raw: string): Error | null {
        const cleaned = (raw ?? "").trim();

        if (cleaned.length === 0) {
            return new Error("Short must not be blank");
        }

        if (cleaned.length > 120) {
            return new Error("Short must be <= 120 chars");
        }

        this._value = cleaned;
        return null;
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
        <label for="${this._name}">Short</label>
        <input
          type="text"
          id="${this._name}"
          name="${this._name}"
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
    readonly _name = "description";
    private _value: string;

    constructor(initialValue: string) {
        this._value = "";
        const error = this.setFromRaw(initialValue);
        if (error) {
            throw error;
        }
    }

    setFromRaw(raw: string): Error | null {
        const cleaned = (raw ?? "").trim();

        if (cleaned.length > 5000) {
            return new Error("Description must be <= 5000 chars");
        }

        this._value = cleaned;
        return null;
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
        <label for="${this._name}">Description</label>
        <textarea id="${this._name}" name="${this._name}">${escapeHtml(this._value)}</textarea>
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
    readonly _name = "due_date";
    private _value: string;

    constructor(initialValue: string) {
        this._value = "";
        const error = this.setFromRaw(initialValue);
        if (error) {
            throw error;
        }
    }

    setFromRaw(raw: string): Error | null {
        const cleaned = (raw ?? "").trim();

        if (cleaned === "") {
            this._value = "";
            return null;
        }

        if (!/^\d{4}-\d{2}-\d{2}$/.test(cleaned)) {
            return new Error("Due Date must be YYYY-MM-DD");
        }

        this._value = cleaned;
        return null;
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
        <label for="${this._name}">Due Date</label>
        <input
          type="date"
          id="${this._name}"
          name="${this._name}"
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
    readonly _name = "cost_of_delay";
    private _value: number;

    constructor(initialValue: number) {
        this._value = 0;
        const error = this.setFromNumber(initialValue);
        if (error) {
            throw error;
        }
    }

    setFromRaw(raw: string): Error | null {
        const cleaned = (raw ?? "").trim();
        const parsed = Number.parseInt(cleaned, 10);

        if (Number.isNaN(parsed)) {
            return new Error("Cost Of Delay must be an integer");
        }

        const error = this.setFromNumber(parsed);
        if (error) {
            return error;
        }
        return null;
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
        <label for="${this._name}">Cost Of Delay</label>
        <input
          type="number"
          id="${this._name}"
          name="${this._name}"
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

    private setFromNumber(value: number): Error | null {
        if (value < -2) {
            return new Error("Cost Of Delay must be >= -2");
        }

        if (value > 2) {
            return new Error("Cost Of Delay must be <= 2");
        }

        this._value = value;
        return null;
    }
}

/**
 * Domain-specific selection for effort sizing.
 */
class TodoEffort {
    readonly _name = "effort";
    private _value: string;
    private _options: string[];

    constructor(initialValue: string) {
        this._value = "";
        this._options = ["mins", "hours", "days", "weeks", "months"];
        const error = this.setFromRaw(initialValue);
        if (error) {
            throw error;
        }
    }

    setFromRaw(raw: string): Error | null {
        const cleaned = (raw ?? "").trim();

        if (!this._options.includes(cleaned)) {
            return new Error(`Invalid effort: ${cleaned}`);
        }

        this._value = cleaned;
        return null;
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
        <label for="${this._name}">Effort</label>
        <select id="${this._name}" name="${this._name}">
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
    private readonly _table_name = "todos";
    private readonly _id: number;

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

    apply(raw: TodoRawInput): TodoValidationResult {
        // Validate against a clone so the original todo stays unchanged on failure.
        const trial = this.clone();
        const errors = [] as TodoValidationError[];

        this.pushFieldError(errors, "short", trial.shortField.setFromRaw(raw.short));
        this.pushFieldError(errors, "description", trial.descriptionField.setFromRaw(raw.description));
        this.pushFieldError(errors, "due_date", trial.dueDateField.setFromRaw(raw.due_date));
        this.pushFieldError(errors, "cost_of_delay", trial.costOfDelayField.setFromRaw(raw.cost_of_delay));
        this.pushFieldError(errors, "effort", trial.effortField.setFromRaw(raw.effort));

        if (errors.length > 0) {
            return { ok: false, errors };
        }

        this.shortField = trial.shortField;
        this.descriptionField = trial.descriptionField;
        this.dueDateField = trial.dueDateField;
        this.costOfDelayField = trial.costOfDelayField;
        this.effortField = trial.effortField;
        return { ok: true, errors: [] };
    }

    patch(raw: TodoPatchInput): TodoValidationResult {
        // Only apply provided values, ignore missing ones.
        const trial = this.clone();
        const errors = [] as TodoValidationError[];

        if (raw.short !== undefined) {
            this.pushFieldError(errors, "short", trial.shortField.setFromRaw(raw.short));
        }

        if (raw.description !== undefined) {
            this.pushFieldError(errors, "description", trial.descriptionField.setFromRaw(raw.description));
        }

        if (raw.due_date !== undefined) {
            this.pushFieldError(errors, "due_date", trial.dueDateField.setFromRaw(raw.due_date));
        }

        if (raw.cost_of_delay !== undefined) {
            this.pushFieldError(errors, "cost_of_delay", trial.costOfDelayField.setFromRaw(raw.cost_of_delay));
        }

        if (raw.effort !== undefined) {
            this.pushFieldError(errors, "effort", trial.effortField.setFromRaw(raw.effort));
        }

        if (errors.length > 0) {
            return { ok: false, errors };
        }

        this.shortField = trial.shortField;
        this.descriptionField = trial.descriptionField;
        this.dueDateField = trial.dueDateField;
        this.costOfDelayField = trial.costOfDelayField;
        this.effortField = trial.effortField;
        return { ok: true, errors: [] };
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

    toJSON(): TodoJson {
        return this.toJson();
    }

    insertSql( ): { sql: string; params: SqlParam[] } {
        const params: SqlParam[] = [
            this.shortField.valueAsSqlParam(),
            this.descriptionField.valueAsSqlParam(),
            this.dueDateField.valueAsSqlParam(),
            this.costOfDelayField.valueAsSqlParam(),
            this.effortField.valueAsSqlParam(),
        ];
        const INSERT = `INSERT INTO ${this._table_name}`; 
        const FIELDS = `(${this.shortField._name}, ${this.descriptionField._name}, ${this.dueDateField._name}, ${this.costOfDelayField._name}, ${this.effortField._name})`; 
        const VALUES = 'VALUES ($1, $2, $3, $4, $5)';  
        return {
            sql: `${INSERT} ${FIELDS} ${VALUES}`,
            params,
        };
    }

    updateSql(): { sql: string; params: SqlParam[] } {
        const params: SqlParam[] = [
            this.shortField.valueAsSqlParam(),
            this.descriptionField.valueAsSqlParam(),
            this.dueDateField.valueAsSqlParam(),
            this.costOfDelayField.valueAsSqlParam(),
            this.effortField.valueAsSqlParam(),
        ];
        const UPDATE = `UPDATE ${this._table_name}`;
        const SET = `SET ${this.shortField._name} = $1, ${this.descriptionField._name} = $2, ${this.dueDateField._name} = $3, ${this.costOfDelayField._name} = $4, ${this.effortField._name} = $5`;
        const WHERE = `WHERE id = $6`;
        return {
            sql: `${UPDATE} ${SET} ${WHERE}`,
            params: [...params, this._id],
        };
    }

    private clone(): Todo {
        return new Todo(this._id, {
            short: this.short(),
            description: this.description(),
            due_date: this.dueDate(),
            cost_of_delay: this.costOfDelay(),
            effort: this.effort(),
        });
    }

    private pushFieldError(errors: TodoValidationError[], field: string, error: Error | null): void {
        if (!error) {
            return;
        }

        errors.push({ field, message: error.message });
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

export { 
    Todo, 
    TodoInitial, 
    TodoRawInput, 
    TodoPatchInput, 
    TodoJson, 
    TodoValidationError, 
    TodoValidationResult 
};