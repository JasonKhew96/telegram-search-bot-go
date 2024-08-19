// Code generated by SQLBoiler 4.16.2 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// Message is an object representing the database table.
type Message struct {
	ID        string    `boil:"id" json:"id" toml:"id" yaml:"id"`
	ChatID    int64     `boil:"chat_id" json:"chat_id" toml:"chat_id" yaml:"chat_id"`
	FromID    int64     `boil:"from_id" json:"from_id" toml:"from_id" yaml:"from_id"`
	MSGID     int64     `boil:"msg_id" json:"msg_id" toml:"msg_id" yaml:"msg_id"`
	Text      string    `boil:"text" json:"text" toml:"text" yaml:"text"`
	Timestamp time.Time `boil:"timestamp" json:"timestamp" toml:"timestamp" yaml:"timestamp"`
	DeletedAt null.Time `boil:"deleted_at" json:"deleted_at,omitempty" toml:"deleted_at" yaml:"deleted_at,omitempty"`

	R *messageR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L messageL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var MessageColumns = struct {
	ID        string
	ChatID    string
	FromID    string
	MSGID     string
	Text      string
	Timestamp string
	DeletedAt string
}{
	ID:        "id",
	ChatID:    "chat_id",
	FromID:    "from_id",
	MSGID:     "msg_id",
	Text:      "text",
	Timestamp: "timestamp",
	DeletedAt: "deleted_at",
}

var MessageTableColumns = struct {
	ID        string
	ChatID    string
	FromID    string
	MSGID     string
	Text      string
	Timestamp string
	DeletedAt string
}{
	ID:        "message.id",
	ChatID:    "message.chat_id",
	FromID:    "message.from_id",
	MSGID:     "message.msg_id",
	Text:      "message.text",
	Timestamp: "message.timestamp",
	DeletedAt: "message.deleted_at",
}

// Generated where

type whereHelpertime_Time struct{ field string }

func (w whereHelpertime_Time) EQ(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.EQ, x)
}
func (w whereHelpertime_Time) NEQ(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.NEQ, x)
}
func (w whereHelpertime_Time) LT(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpertime_Time) LTE(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpertime_Time) GT(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpertime_Time) GTE(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

type whereHelpernull_Time struct{ field string }

func (w whereHelpernull_Time) EQ(x null.Time) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpernull_Time) NEQ(x null.Time) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpernull_Time) LT(x null.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpernull_Time) LTE(x null.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpernull_Time) GT(x null.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpernull_Time) GTE(x null.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

func (w whereHelpernull_Time) IsNull() qm.QueryMod    { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpernull_Time) IsNotNull() qm.QueryMod { return qmhelper.WhereIsNotNull(w.field) }

var MessageWhere = struct {
	ID        whereHelperstring
	ChatID    whereHelperint64
	FromID    whereHelperint64
	MSGID     whereHelperint64
	Text      whereHelperstring
	Timestamp whereHelpertime_Time
	DeletedAt whereHelpernull_Time
}{
	ID:        whereHelperstring{field: "\"message\".\"id\""},
	ChatID:    whereHelperint64{field: "\"message\".\"chat_id\""},
	FromID:    whereHelperint64{field: "\"message\".\"from_id\""},
	MSGID:     whereHelperint64{field: "\"message\".\"msg_id\""},
	Text:      whereHelperstring{field: "\"message\".\"text\""},
	Timestamp: whereHelpertime_Time{field: "\"message\".\"timestamp\""},
	DeletedAt: whereHelpernull_Time{field: "\"message\".\"deleted_at\""},
}

// MessageRels is where relationship names are stored.
var MessageRels = struct {
}{}

// messageR is where relationships are stored.
type messageR struct {
}

// NewStruct creates a new relationship struct
func (*messageR) NewStruct() *messageR {
	return &messageR{}
}

// messageL is where Load methods for each relationship are stored.
type messageL struct{}

var (
	messageAllColumns            = []string{"id", "chat_id", "from_id", "msg_id", "text", "timestamp", "deleted_at"}
	messageColumnsWithoutDefault = []string{"id", "chat_id", "from_id", "msg_id", "text", "timestamp"}
	messageColumnsWithDefault    = []string{"deleted_at"}
	messagePrimaryKeyColumns     = []string{"id"}
	messageGeneratedColumns      = []string{}
)

type (
	// MessageSlice is an alias for a slice of pointers to Message.
	// This should almost always be used instead of []Message.
	MessageSlice []*Message
	// MessageHook is the signature for custom Message hook methods
	MessageHook func(context.Context, boil.ContextExecutor, *Message) error

	messageQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	messageType                 = reflect.TypeOf(&Message{})
	messageMapping              = queries.MakeStructMapping(messageType)
	messagePrimaryKeyMapping, _ = queries.BindMapping(messageType, messageMapping, messagePrimaryKeyColumns)
	messageInsertCacheMut       sync.RWMutex
	messageInsertCache          = make(map[string]insertCache)
	messageUpdateCacheMut       sync.RWMutex
	messageUpdateCache          = make(map[string]updateCache)
	messageUpsertCacheMut       sync.RWMutex
	messageUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var messageAfterSelectMu sync.Mutex
var messageAfterSelectHooks []MessageHook

var messageBeforeInsertMu sync.Mutex
var messageBeforeInsertHooks []MessageHook
var messageAfterInsertMu sync.Mutex
var messageAfterInsertHooks []MessageHook

var messageBeforeUpdateMu sync.Mutex
var messageBeforeUpdateHooks []MessageHook
var messageAfterUpdateMu sync.Mutex
var messageAfterUpdateHooks []MessageHook

var messageBeforeDeleteMu sync.Mutex
var messageBeforeDeleteHooks []MessageHook
var messageAfterDeleteMu sync.Mutex
var messageAfterDeleteHooks []MessageHook

var messageBeforeUpsertMu sync.Mutex
var messageBeforeUpsertHooks []MessageHook
var messageAfterUpsertMu sync.Mutex
var messageAfterUpsertHooks []MessageHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Message) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range messageAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Message) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range messageBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Message) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range messageAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Message) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range messageBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Message) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range messageAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Message) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range messageBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Message) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range messageAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Message) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range messageBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Message) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range messageAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddMessageHook registers your hook function for all future operations.
func AddMessageHook(hookPoint boil.HookPoint, messageHook MessageHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		messageAfterSelectMu.Lock()
		messageAfterSelectHooks = append(messageAfterSelectHooks, messageHook)
		messageAfterSelectMu.Unlock()
	case boil.BeforeInsertHook:
		messageBeforeInsertMu.Lock()
		messageBeforeInsertHooks = append(messageBeforeInsertHooks, messageHook)
		messageBeforeInsertMu.Unlock()
	case boil.AfterInsertHook:
		messageAfterInsertMu.Lock()
		messageAfterInsertHooks = append(messageAfterInsertHooks, messageHook)
		messageAfterInsertMu.Unlock()
	case boil.BeforeUpdateHook:
		messageBeforeUpdateMu.Lock()
		messageBeforeUpdateHooks = append(messageBeforeUpdateHooks, messageHook)
		messageBeforeUpdateMu.Unlock()
	case boil.AfterUpdateHook:
		messageAfterUpdateMu.Lock()
		messageAfterUpdateHooks = append(messageAfterUpdateHooks, messageHook)
		messageAfterUpdateMu.Unlock()
	case boil.BeforeDeleteHook:
		messageBeforeDeleteMu.Lock()
		messageBeforeDeleteHooks = append(messageBeforeDeleteHooks, messageHook)
		messageBeforeDeleteMu.Unlock()
	case boil.AfterDeleteHook:
		messageAfterDeleteMu.Lock()
		messageAfterDeleteHooks = append(messageAfterDeleteHooks, messageHook)
		messageAfterDeleteMu.Unlock()
	case boil.BeforeUpsertHook:
		messageBeforeUpsertMu.Lock()
		messageBeforeUpsertHooks = append(messageBeforeUpsertHooks, messageHook)
		messageBeforeUpsertMu.Unlock()
	case boil.AfterUpsertHook:
		messageAfterUpsertMu.Lock()
		messageAfterUpsertHooks = append(messageAfterUpsertHooks, messageHook)
		messageAfterUpsertMu.Unlock()
	}
}

// One returns a single message record from the query.
func (q messageQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Message, error) {
	o := &Message{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for message")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all Message records from the query.
func (q messageQuery) All(ctx context.Context, exec boil.ContextExecutor) (MessageSlice, error) {
	var o []*Message

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Message slice")
	}

	if len(messageAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all Message records in the query.
func (q messageQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count message rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q messageQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if message exists")
	}

	return count > 0, nil
}

// Messages retrieves all the records using an executor.
func Messages(mods ...qm.QueryMod) messageQuery {
	mods = append(mods, qm.From("\"message\""), qmhelper.WhereIsNull("\"message\".\"deleted_at\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"message\".*"})
	}

	return messageQuery{q}
}

// FindMessage retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindMessage(ctx context.Context, exec boil.ContextExecutor, iD string, selectCols ...string) (*Message, error) {
	messageObj := &Message{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"message\" where \"id\"=? and \"deleted_at\" is null", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, messageObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from message")
	}

	if err = messageObj.doAfterSelectHooks(ctx, exec); err != nil {
		return messageObj, err
	}

	return messageObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Message) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no message provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(messageColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	messageInsertCacheMut.RLock()
	cache, cached := messageInsertCache[key]
	messageInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			messageAllColumns,
			messageColumnsWithDefault,
			messageColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(messageType, messageMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(messageType, messageMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"message\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"message\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into message")
	}

	if !cached {
		messageInsertCacheMut.Lock()
		messageInsertCache[key] = cache
		messageInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the Message.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Message) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	messageUpdateCacheMut.RLock()
	cache, cached := messageUpdateCache[key]
	messageUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			messageAllColumns,
			messagePrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update message, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"message\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 0, wl),
			strmangle.WhereClause("\"", "\"", 0, messagePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(messageType, messageMapping, append(wl, messagePrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update message row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for message")
	}

	if !cached {
		messageUpdateCacheMut.Lock()
		messageUpdateCache[key] = cache
		messageUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q messageQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for message")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for message")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o MessageSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), messagePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"message\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, messagePrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in message slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all message")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Message) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no message provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(messageColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	messageUpsertCacheMut.RLock()
	cache, cached := messageUpsertCache[key]
	messageUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			messageAllColumns,
			messageColumnsWithDefault,
			messageColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			messageAllColumns,
			messagePrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert message, could not build update column list")
		}

		ret := strmangle.SetComplement(messageAllColumns, strmangle.SetIntersect(insert, update))

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(messagePrimaryKeyColumns))
			copy(conflict, messagePrimaryKeyColumns)
		}
		cache.query = buildUpsertQuerySQLite(dialect, "\"message\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(messageType, messageMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(messageType, messageMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if errors.Is(err, sql.ErrNoRows) {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert message")
	}

	if !cached {
		messageUpsertCacheMut.Lock()
		messageUpsertCache[key] = cache
		messageUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single Message record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Message) Delete(ctx context.Context, exec boil.ContextExecutor, hardDelete bool) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Message provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	var (
		sql  string
		args []interface{}
	)
	if hardDelete {
		args = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), messagePrimaryKeyMapping)
		sql = "DELETE FROM \"message\" WHERE \"id\"=?"
	} else {
		currTime := time.Now().In(boil.GetLocation())
		o.DeletedAt = null.TimeFrom(currTime)
		wl := []string{"deleted_at"}
		sql = fmt.Sprintf("UPDATE \"message\" SET %s WHERE \"id\"=?",
			strmangle.SetParamNames("\"", "\"", 0, wl),
		)
		valueMapping, err := queries.BindMapping(messageType, messageMapping, append(wl, messagePrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
		args = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), valueMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from message")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for message")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q messageQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor, hardDelete bool) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no messageQuery provided for delete all")
	}

	if hardDelete {
		queries.SetDelete(q.Query)
	} else {
		currTime := time.Now().In(boil.GetLocation())
		queries.SetUpdate(q.Query, M{"deleted_at": currTime})
	}

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from message")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for message")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o MessageSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor, hardDelete bool) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(messageBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var (
		sql  string
		args []interface{}
	)
	if hardDelete {
		for _, obj := range o {
			pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), messagePrimaryKeyMapping)
			args = append(args, pkeyArgs...)
		}
		sql = "DELETE FROM \"message\" WHERE " +
			strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, messagePrimaryKeyColumns, len(o))
	} else {
		currTime := time.Now().In(boil.GetLocation())
		for _, obj := range o {
			pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), messagePrimaryKeyMapping)
			args = append(args, pkeyArgs...)
			obj.DeletedAt = null.TimeFrom(currTime)
		}
		wl := []string{"deleted_at"}
		sql = fmt.Sprintf("UPDATE \"message\" SET %s WHERE "+
			strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, messagePrimaryKeyColumns, len(o)),
			strmangle.SetParamNames("\"", "\"", 0, wl),
		)
		args = append([]interface{}{currTime}, args...)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from message slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for message")
	}

	if len(messageAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Message) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindMessage(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *MessageSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := MessageSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), messagePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"message\".* FROM \"message\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, messagePrimaryKeyColumns, len(*o)) +
		"and \"deleted_at\" is null"

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in MessageSlice")
	}

	*o = slice

	return nil
}

// MessageExists checks if the Message row exists.
func MessageExists(ctx context.Context, exec boil.ContextExecutor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"message\" where \"id\"=? and \"deleted_at\" is null limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if message exists")
	}

	return exists, nil
}

// Exists checks if the Message row exists.
func (o *Message) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return MessageExists(ctx, exec, o.ID)
}
