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

type Short = string;
type Description = string;
type DueDate = string
type CostOfDelay = -2 | -1 | 0 | 1 | 2
type Effort = 'mins' | 'hours' | 'days' | 'weeks' | 'months';

type TodoRawInput = {
    short: Short;
    description: Description;
    due_date: DueDate;
    cost_of_delay: CostOfDelay;
    effort: Effort;
};
type TodoOptionals = {
    short?: Short;
    description?: Description;
    due_date?: DueDate;
    cost_of_delay?: CostOfDelay;
    effort?: Effort;
};
type TodoPatchInput = TodoOptionals;
type TodoInitial = TodoOptionals;

interface TodoField {
    _name: string;
    value(): Short | Description | DueDate | CostOfDelay | Effort;
    setFromRaw(raw: Short | Description | DueDate | CostOfDelay | Effort): Error | null;
    renderField(): string;
}

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
    short: Short;
    description: Description;
    due_date: DueDate | null;
    cost_of_delay: CostOfDelay;
    effort: Effort;
};

// HELPERS:
// HELPERS/SQL:

/**
 * Dedicated field for the todo title. Required and intentionally strict.
 */
class TodoShort implements TodoField {
    readonly _name = 'short';
    private _value: Short;

    constructor(initialValue: Short) {
        this._value = "";
        const error = this.setFromRaw(initialValue);
        if (error) {
            throw error;
        }
    }

    setFromRaw(raw: Short): Error | null {
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

    value(): Short {
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


}

/**
 * Dedicated field for detailed description. Optional but bounded.
 */
class TodoDescription implements TodoField {
    readonly _name = "description";
    private _value: Description;

    constructor(initialValue: Description) {
        this._value = "";
        const error = this.setFromRaw(initialValue);
        if (error) {
            throw error;
        }
    }

    setFromRaw(raw: Description): Error | null {
        const cleaned = (raw ?? "").trim();

        if (cleaned.length > 5000) {
            return new Error("Description must be <= 5000 chars");
        }

        this._value = cleaned;
        return null;
    }

    value(): Description {
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


}

/**
 * Dedicated due date with a narrow accepted format for consistency with SQL and forms.
 */
class TodoDueDate implements TodoField {
    readonly _name = "due_date";
    private _value: DueDate;

    constructor(initialValue: DueDate) {
        this._value = "";
        const error = this.setFromRaw(initialValue);
        if (error) {
            throw error;
        }
    }

    setFromRaw(raw: DueDate): Error | null {
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

    value(): DueDate {
        return this._value;
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

}


/**
 * Domain-specific integer: cost of delay for this todo, constrained to a small scale.
 */
class TodoCostOfDelay implements TodoField {
    readonly _name = "cost_of_delay";
    private _value: CostOfDelay;

    constructor(initialValue: CostOfDelay) {
        this._value = 0;
        const error = this.setFromNumber(initialValue);
        if (error) {
            throw error;
        }
    }

    setFromRaw(raw: CostOfDelay): Error | null {


        if (Number.isNaN(raw)) {
            return new Error("Cost Of Delay must be an integer");
        }

        const error = this.setFromNumber(raw);
        if (error) {
            return error;
        }
        return null;
    }

    value(): CostOfDelay {
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

    private setFromNumber(value: CostOfDelay): Error | null {
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
class TodoEffort implements TodoField {
    readonly _name = "effort";
    private _value: Effort;
    private _options: Effort[];

    constructor(initialValue: Effort) {
        this._value = "hours";
        this._options = ["mins", "hours", "days", "weeks", "months"];
        const error = this.setFromRaw(initialValue);
        if (error) {
            throw error;
        }
    }

    setFromRaw(raw: Effort): Error | null {
        if (!this._options.includes(raw)) {
            return new Error(`Invalid effort: ${raw}`);
        }

        this._value = raw;
        return null;
    }

    value(): Effort {
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

}

/**
 * The Todo class. Some would call it a god-object, but I view it as a deep module (for now). I might changem my view.
 * 
 * The Todo class should enable users to create Todo items and change them safely. It should enable them to persist them without issues and to render representations and forms in an easy manner.
 */
class Todo {
    readonly _table_name = "todos";
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

    private fields(): TodoField[] {
        return [
            this.shortField,
            this.descriptionField,
            this.dueDateField,
            this.costOfDelayField,
            this.effortField,
        ];
    }

    id(): number {
        return this._id;
    }
    /**
     * A short title for the todo
     */
    short(): string {
        return this.shortField.value();
    }

    /**
     * A detailed description of the todo. Might be empty
     */
    description(): string {
        return this.descriptionField.value();
    }

    /**
     * The date the task ideally should be done. Might be empty
     */
    dueDate(): string {
        return this.dueDateField.value();
    }

    /**
     * A raw estimate how much we loose if we postpone this task, on a scale from -2 (allmost nothing) to 2 (a ton). Required.
     */
    costOfDelay(): number {
        return Number(this.costOfDelayField.value());
    }

    /**
     * A raw estimate of effort in time. Allowed are "mins", "hours", "days", "weeks" and "months". Required.
     */
    effort(): string {
        return this.effortField.value();
    }

    /**
     * Update all fields. Make sure to include all data and valid data or this will not update and return errors. Use patch() if you want to update only some fields and ignore missing ones.
     */
    apply(raw: TodoRawInput): TodoValidationResult {
        const errors = [] as TodoValidationError[];
        // Validate against a clone so the original todo stays unchanged on failure.
        const trial = this.clone();
        const fields = trial.fields();
        for (const field of fields) {
            const key = field._name as keyof TodoRawInput;
            const error = field.setFromRaw(raw[key]);
            this.pushFieldError(errors, field._name, error);
        }

        if (errors.length > 0) {
            return { ok: false, errors };
        }

        for (const field of fields) {
            const key = field._name as keyof TodoRawInput;
            field.setFromRaw(raw[key]);
        }

        return { ok: true, errors: [] };
    }

    /**
     * Updates a partial set of fields. But they need to be valid or the update fails and returns errors.
     */
    patch(raw: TodoPatchInput): TodoValidationResult {
        const errors = [] as TodoValidationError[];
        // Only apply provided values, ignore missing ones.
        const trial = this.clone();

        const fields = trial.fields();

        for (const field of fields) {
            const key = field._name as keyof TodoPatchInput;

            if (raw[key] !== undefined) {
                const error = field.setFromRaw(raw[key]);
                trial.pushFieldError(errors, field._name, error);
            }
        }

        if (errors.length > 0) {
            return { ok: false, errors };
        }


        for (const field of fields) {
            const key = field._name as keyof TodoPatchInput;
            if (raw[key] !== undefined) {
                field.setFromRaw(raw[key] as string);
            }
        }
        return { ok: true, errors: [] };
    }

    /**
     * You get a table row for a quick overview
     */
    renderTableRow(): string {
        return `
      <tr>
        <td>${this._id}</td>
        <td>${escapeHtml(this.shortField.value())}</td>
        <td>${escapeHtml(this.descriptionField.value())}</td>
        <td>${escapeHtml(this.dueDateField.value())}</td>
        <td>${escapeHtml(String(this.costOfDelayField.value()))}</td>
        <td>${escapeHtml(this.effortField.value())}</td>
      </tr>
    `.trim();
    }

    /**
     * You get a card with all details for popups on small devices or similar use cases.
     */
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

    /**
     * This gives you a form where you can edit all fields. You can use it for both creating and editing. Just make sure to provide the right action URL. The form is unstyled and basic on purpose, but it includes all necessary attributes for a good user experience like labels, input types and max lengths.
     */
    renderEditForm(action: string): string {
        const fieldsHtml = this.fields().map((field) => field.renderField()).join("\n");

        return `
      <form method="post" action="${escapeHtml(action)}">
        ${fieldsHtml}
        <div>
          <input type="submit" value="Save">
        </div>
      </form>
    `.trim();
    }

    /**
     * A custom JSON representation close to the domain model. This is what we use for persistence and also for API responses. We keep it close to the model so we have a single source of truth for how the data looks like in JSON. If we wanted to add a different representation for a specific frontend, we could add another method like toApiJson() or similar.
     */
    toJson(): TodoJson {
        return {
            id: this._id,
            short: this.shortField.value(),
            description: this.descriptionField.value(),
            due_date: this.dueDateField.value() || null,
            cost_of_delay: this.costOfDelayField.value(),
            effort: this.effortField.value(),
        };
    }

    /**
     * An automatic converson to JSON from the Object itself. I have no idea for what it's useful..
     */
    toJSON(): TodoJson {
        return this.toJson();
    }

    /**
     * You get all keys of the todo as an array. This is useful for generic handling of fields, for example in forms or similar use cases. It also abstracts away the internal structure of the Todo class so the consumer doesn't need to know about the specific field classes we use internally.
     */
    keys(): string[] {
        return this.fields().map((field) => field._name);
    }

    /** Gives you a very plain object for use in query building (e.g SQL) */
    values(): Record<string, string> {
        const object = {
            short: this.shortField.value(),
            description: this.descriptionField.value(),
            due_date: this.dueDateField.value(),
            cost_of_delay: String(this.costOfDelayField.value()),
            effort: this.effortField.value(),
        } as Record<string, string>;
        return object;
    }
 
    private clone(): Todo {
        return new Todo(
            this._id, {
            short: this.shortField.value(),
            description: this.descriptionField.value(),
            due_date: this.dueDateField.value(),
            cost_of_delay: this.costOfDelayField.value(),
            effort: this.effortField.value(),
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
    type TodoInitial,
    type TodoJson,
    type TodoRawInput,
    type TodoPatchInput,
    type TodoValidationResult,
    type TodoValidationError,    
};