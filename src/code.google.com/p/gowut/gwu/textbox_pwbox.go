// Copyright 2013 Andras Belicza. All rights reserved.

// Defines the TextBox component.

package gwu

import (
	"net/http"
	"strconv"
)

// TextBox interface defines a component for text input purpose.
// 
// Suggested event type to handle actions: ETYPE_CHANGE
//
// By default the value of the TextBox is synchronized with the server
// on ETYPE_CHANGE event which is when the TextBox loses focus
// or when the ENTER key is pressed.
// If you want a TextBox to synchronize values during editing
// (while you type in characters), add the ETYPE_KEY_UP event type
// to the events on which synchronization happens by calling:
// 		AddSyncOnETypes(ETYPE_KEY_UP)
// 
// Default style class: "gwu-TextBox"
type TextBox interface {
	// TextBox is a component.
	Comp

	// TextBox has text.
	HasText

	// TextBox can be enabled/disabled.
	HasEnabled

	// ReadOnly returns if the text box is read-only.
	ReadOnly() bool

	// SetReadOnly sets if the text box is read-only.
	SetReadOnly(readOnly bool)

	// Rows returns the number of displayed rows.
	Rows() int

	// SetRows sets the number of displayed rows.
	// rows=1 will make this a simple, one-line input text box,
	// rows>1 will make this a text area.
	SetRows(rows int)

	// Cols returns the number of displayed columns.
	Cols() int

	// SetCols sets the number of displayed columns.
	SetCols(cols int)

	// MaxLength returns the maximum number of characters
	// allowed in the text box.
	// -1 is returned if there is no maximum length set.
	MaxLength() int

	// SetMaxLength sets the maximum number of characters
	// allowed in the text box.
	// Pass -1 to not limit the maximum length.
	SetMaxLength(maxLength int)
}

// PasswBox interface defines a text box for password input purpose.
// 
// Suggested event type to handle actions: ETYPE_CHANGE
//
// By default the value of the PasswBox is synchronized with the server
// on ETYPE_CHANGE event which is when the PasswBox loses focus
// or when the ENTER key is pressed.
// If you want a PasswBox to synchronize values during editing
// (while you type in characters), add the ETYPE_KEY_UP event type
// to the events on which synchronization happens by calling:
// 		AddSyncOnETypes(ETYPE_KEY_UP)
// 
// Default style class: "gwu-PasswBox"
type PasswBox interface {
	// PasswBox is a TextBox.
	TextBox
}

// TextBox implementation.
type textBoxImpl struct {
	compImpl       // Component implementation
	hasTextImpl    // Has text implementation
	hasEnabledImpl // Has enabled implementation

	isPassw    bool // Tells if the text box is a password box
	rows, cols int  // Number of displayed rows and columns.
}

// NewTextBox creates a new TextBox.
func NewTextBox(text string) TextBox {
	c := newTextBoxImpl("encodeURIComponent(this.value)", text, false)
	c.Style().AddClass("gwu-TextBox")
	return &c
}

// NewPasswBox creates a new PasswBox.
func NewPasswBox(text string) TextBox {
	c := newTextBoxImpl("encodeURIComponent(this.value)", text, true)
	c.Style().AddClass("gwu-PasswBox")
	return &c
}

// newTextBoxImpl creates a new textBoxImpl.
func newTextBoxImpl(valueProviderJs, text string, isPassw bool) textBoxImpl {
	c := textBoxImpl{newCompImpl(valueProviderJs), newHasTextImpl(text), newHasEnabledImpl(), isPassw, 1, 20}
	c.AddSyncOnETypes(ETYPE_CHANGE)
	return c
}

func (c *textBoxImpl) ReadOnly() bool {
	ro := c.Attr("readonly")
	return len(ro) > 0
}

func (c *textBoxImpl) SetReadOnly(readOnly bool) {
	if readOnly {
		c.SetAttr("readonly", "readonly")
	} else {
		c.SetAttr("readonly", "")
	}
}

func (c *textBoxImpl) Rows() int {
	return c.rows
}

func (c *textBoxImpl) SetRows(rows int) {
	c.rows = rows
}

func (c *textBoxImpl) Cols() int {
	return c.cols
}

func (c *textBoxImpl) SetCols(cols int) {
	c.cols = cols
}

func (c *textBoxImpl) MaxLength() int {
	if ml := c.Attr("maxlength"); len(ml) > 0 {
		if i, err := strconv.Atoi(ml); err == nil {
			return i
		}
	}
	return -1
}

func (c *textBoxImpl) SetMaxLength(maxLength int) {
	if maxLength < 0 {
		c.SetAttr("maxlength", "")
	} else {
		c.SetAttr("maxlength", strconv.Itoa(maxLength))
	}
}

func (c *textBoxImpl) preprocessEvent(event Event, r *http.Request) {
	// Empty string for text box is a valid value.
	// So we have to check whether it is supplied, not just whether its len() > 0 
	value := r.FormValue(_PARAM_COMP_VALUE)
	if len(value) > 0 {
		c.text = value
	} else {
		// Empty string might be a valid value, if the component value param is present:
		values, present := r.Form[_PARAM_COMP_VALUE] // Form is surely parsed (we called FormValue())
		if present && len(values) > 0 {
			c.text = values[0]
		}
	}
}

func (c *textBoxImpl) Render(w writer) {
	if c.rows <= 1 || c.isPassw {
		c.renderInput(w)
	} else {
		c.renderTextArea(w)
	}
}

// renderInput renders the component as an input HTML tag.
func (c *textBoxImpl) renderInput(w writer) {
	w.Writes("<input type=\"")
	if c.isPassw {
		w.Writes("password")
	} else {
		w.Writes("text")
	}
	w.Writevs("\" size=\"", c.cols)
	w.Write(_STR_QUOTE)
	c.renderAttrsAndStyle(w)
	c.renderEnabled(w)
	c.renderEHandlers(w)

	w.Writes(" value=\"")
	c.renderText(w)
	w.Writes("\"/>")
}

// renderTextArea renders the component as an textarea HTML tag.
func (c *textBoxImpl) renderTextArea(w writer) {
	w.Writes("<textarea")
	c.renderAttrsAndStyle(w)
	c.renderEnabled(w)
	c.renderEHandlers(w)

	// New line char after the <textarea> tag is ignored.
	// So we must render a newline after textarea, else if text value
	// starts with a new line, it will be ommitted!
	w.Writevs(" rows=\"", c.rows, "\" cols=\"", c.cols, "\">\n")

	c.renderText(w)
	w.Writes("</textarea>")
}
