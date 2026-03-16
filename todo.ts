interface Field {
    name(): string;
    label(): string;

    setFromRaw(raw: string): void;

    valueAsString(): string;
    valueAsSqlParam(): string | number | null;
    valueAsJson(): string | number | null;

    renderField(): string;
    renderTableCell(): string;
}

class TextField implements Field {
    private _value: string;
    private _name: string;
    private _label?: string;
    private _required: boolean;
    private _maxLength: number;
    private _multiline: boolean;

    constructor(
        name: string,
        initialValue: string,
        required: boolean,
        maxLength: number,
        multiline = false,
        label?: string
    ) {
        this._name = name;
        this._label = label;
        this._required = required;
        this._maxLength = maxLength;
        this._multiline = multiline;
        this._value = "";
        this.setFromRaw(initialValue);
    }

    name(): string {
        return this._name;
    }

    label(): string {
        return this._label ?? humanise(this._name);
    }

    setFromRaw(raw: string): void {
        const cleaned = (raw ?? "").trim();

        if (this._required && cleaned.length === 0) {
            throw new Error(`${this.label()} must not be blank`);
        }

        if (cleaned.length > this._maxLength) {
            throw new Error(`${this.label()} must be <= ${this._maxLength} chars`);
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
        const label = `<label for="${escapeHtml(this._name)}">${escapeHtml(this.label())}</label>`;

        if (this._multiline) {
            return `
        <div>
          ${label}
          <textarea id="${escapeHtml(this._name)}" name="${escapeHtml(this._name)}">${escapeHtml(this._value)}</textarea>
        </div>
      `.trim();
        }

        return `
      <div>
        ${label}
        <input
          type="text"
          id="${escapeHtml(this._name)}"
          name="${escapeHtml(this._name)}"
          value="${escapeHtml(this._value)}"
        >
      </div>
    `.trim();
    }

    renderTableCell(): string {
        return `<td>${escapeHtml(this._value)}</td>`;
    }
}

class DateField implements Field {
    private _value: string;
    private _name: string;
    private _label?: string;

    constructor(name: string, initialValue = "", label?: string) {
        this._name = name;
        this._label = label;
        this._value = "";
        this.setFromRaw(initialValue);
    }

    name(): string {
        return this._name;
    }

    label(): string {
        return this._label ?? humanise(this._name);
    }

    setFromRaw(raw: string): void {
        const cleaned = (raw ?? "").trim();

        if (cleaned === "") {
            this._value = "";
            return;
        }

        if (!/^\d{4}-\d{2}-\d{2}$/.test(cleaned)) {
            throw new Error(`${this.label()} must be YYYY-MM-DD`);
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
        <label for="${escapeHtml(this._name)}">${escapeHtml(this.label())}</label>
        <input
          type="date"
          id="${escapeHtml(this._name)}"
          name="${escapeHtml(this._name)}"
          value="${escapeHtml(this._value)}"
        >
      </div>
    `.trim();
    }

    renderTableCell(): string {
        return `<td>${escapeHtml(this._value)}</td>`;
    }
}

class IntField implements Field {
    private _value: number;
    private _name: string;
    private _label?: string;
    private _min: number;
    private _max: number;

    constructor(name: string, initialValue: number, min: number, max: number, label?: string) {
        this._name = name;
        this._label = label;
        this._min = min;
        this._max = max;
        this._value = initialValue;
        this.validate(initialValue);
    }

    name(): string {
        return this._name;
    }

    label(): string {
        return this._label ?? humanise(this._name);
    }

    setFromRaw(raw: string): void {
        const cleaned = (raw ?? "").trim();
        const parsed = Number.parseInt(cleaned, 10);

        if (Number.isNaN(parsed)) {
            throw new Error(`${this.label()} must be an integer`);
        }

        this.validate(parsed);
        this._value = parsed;
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
        <label for="${escapeHtml(this._name)}">${escapeHtml(this.label())}</label>
        <input
          type="number"
          id="${escapeHtml(this._name)}"
          name="${escapeHtml(this._name)}"
          value="${this._value}"
          min="${this._min}"
          max="${this._max}"
        >
      </div>
    `.trim();
    }

    renderTableCell(): string {
        return `<td>${this._value}</td>`;
    }

    private validate(value: number): void {
        if (value < this._min) {
            throw new Error(`${this.label()} must be >= ${this._min}`);
        }

        if (value > this._max) {
            throw new Error(`${this.label()} must be <= ${this._max}`);
        }
    }
}

class SelectField implements Field {
    private _value: string;
    private _name: string;
    private _label?: string;
    private _options: string[];

    constructor(name: string, options: string[], initialValue: string, label?: string) {
        this._name = name;
        this._label = label;
        this._options = options;
        this._value = "";
        this.setFromRaw(initialValue);
    }

    name(): string {
        return this._name;
    }

    label(): string {
        return this._label ?? humanise(this._name);
    }

    setFromRaw(raw: string): void {
        const cleaned = (raw ?? "").trim();

        if (!this._options.includes(cleaned)) {
            throw new Error(`Invalid value for ${this.label()}: ${cleaned}`);
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
        <label for="${escapeHtml(this._name)}">${escapeHtml(this.label())}</label>
        <select id="${escapeHtml(this._name)}" name="${escapeHtml(this._name)}">
          ${optionsHtml}
        </select>
      </div>
    `.trim();
    }

    renderTableCell(): string {
        return `<td>${escapeHtml(this._value)}</td>`;
    }
}

class Todo {
    private _id: number;

    private shortField: TextField;
    private descriptionField: TextField;
    private dueDateField: DateField;
    private costOfDelayField: IntField;
    private effortField: SelectField;

    constructor(
        id: number,
        initial?: {
            short?: string;
            description?: string;
            due_date?: string;
            cost_of_delay?: number;
            effort?: string;
        }
    ) {
        this._id = id;

        this.shortField = new TextField(
            "short",
            initial?.short ?? "",
            true,
            120
        );

        this.descriptionField = new TextField(
            "description",
            initial?.description ?? "",
            false,
            5000,
            true
        );

        this.dueDateField = new DateField(
            "due_date",
            initial?.due_date ?? ""
        );

        this.costOfDelayField = new IntField(
            "cost_of_delay",
            initial?.cost_of_delay ?? 0,
            -2,
            2
        );

        this.effortField = new SelectField(
            "effort",
            ["mins", "hours", "days", "weeks", "months"],
            initial?.effort ?? "hours"
        );
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

    apply(raw: Record<string, string>): void {
        for (const field of this.fields()) {
            field.setFromRaw(raw[field.name()] ?? "");
        }
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

    toJson(): Record<string, string | number | null> {
        return {
            id: this._id,
            short: this.shortField.valueAsJson(),
            description: this.descriptionField.valueAsJson(),
            due_date: this.dueDateField.valueAsJson(),
            cost_of_delay: this.costOfDelayField.valueAsJson(),
            effort: this.effortField.valueAsJson(),
        };
    }

    insertSql(tableName = "todos"): { sql: string; params: Array<string | number | null> } {
        const fields = this.fields();
        const columns = fields.map((field) => field.name());
        const placeholders = fields.map((_, index) => `$${index + 1}`);
        const params = fields.map((field) => field.valueAsSqlParam());

        return {
            sql: `INSERT INTO ${tableName} (${columns.join(", ")}) VALUES (${placeholders.join(", ")})`,
            params,
        };
    }

    updateSql(tableName = "todos"): { sql: string; params: Array<string | number | null> } {
        const fields = this.fields();
        const assignments = fields.map((field, index) => `${field.name()} = $${index + 1}`);
        const params = fields.map((field) => field.valueAsSqlParam());

        return {
            sql: `UPDATE ${tableName} SET ${assignments.join(", ")} WHERE id = $${fields.length + 1}`,
            params: [...params, this._id],
        };
    }

    private fields(): Field[] {
        return [
            this.shortField,
            this.descriptionField,
            this.dueDateField,
            this.costOfDelayField,
            this.effortField,
        ];
    }
}

function humanise(name: string): string {
    return name
        .split("_")
        .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
        .join(" ");
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