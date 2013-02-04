// Copyright 2013 Andras Belicza. All rights reserved.

// TabPanel component interface and implementation.

package gwu

// TabBar interface defines the tab bar for selecting the visible
// component of a TabPanel.
// 
// Note: Removing a tab component through the tab bar also
// removes the content component from the tab panel of the tab bar.
//
// Default style classes: "gwu-TabBar", "gwu-TabBar-Top", "gwu-TabBar-Bottom",
// "gwu-TabBar-Left", "gwu-TabBar-Right", "gwu-TabBar-NotSelected",
// "gwu-TabBar-Selected"
type TabBar interface {
	// TabBar is a PanelView.
	PanelView
}

// TabBar implementation.
type tabBarImpl struct {
	panelImpl // panel implementation
}

// newTabBarImpl creates a new tabBarImpl.
func newTabBarImpl() *tabBarImpl {
	c := &tabBarImpl{newPanelImpl()}
	c.SetAttr("cellspacing", "0")
	c.SetAttr("cellpadding", "0")
	return c
}

func (c *tabBarImpl) Remove(c2 Comp) bool {
	i := c.CompIdx(c2)
	if i < 0 {
		return false
	}

	// Removing a tab component also needs removing the
	// associated content component. Call parent's (TabPanel) Remove()
	// method which takes care of everything:
	if parent := c.parent; parent != nil {
		if tabPanel, isTabPanel := parent.(TabPanel); isTabPanel {
			return tabPanel.Remove(tabPanel.CompAt(i))
		}
	}

	return c.panelImpl.Remove(c2)
}

// Tab bar placement type.
type TabBarPlacement int

// Tab bar placements.
const (
	TB_PLACEMENT_TOP    TabBarPlacement = iota // Tab bar placement to Top
	TB_PLACEMENT_BOTTOM                        // Tab bar placement to Bottom
	TB_PLACEMENT_LEFT                          // Tab bar placement to Left
	TB_PLACEMENT_RIGHT                         // Tab bar placement to Right
)

// TabPanel interface defines a PanelView which has multiple child components
// but only one is visible at a time. The visible child can be visually selected
// using an internal TabBar component.
// 
// Both the tab panel and its internal tab bar component are PanelViews.
// This gives high layout configuration possibilities.
// Usually you only need to set the tab bar placement with the SetTabBarPlacement()
// method which also sets other reasonable internal layout defaults.
// But you have many other options to override the layout settings.
// If the content component is bigger than the tab bar, you can set the tab bar
// horizontal and the vertical alignment, e.g. with the TabBar().SetAlignment() method.
// To apply cell formatting to individual content components, you can simply use the
// CellFmt() method. If the tab bar is bigger than the content component, you can set
// the content alignment, e.g. with the SetAlignment() method. You can also set different
// alignments for individual tab components using TabBar().CellFmt(). You can also set
// other cell formatting applied to the tab bar using TabBarFmt() method.
// 
// Default style classes: "gwu-TabPanel", "gwu-TabPanel-Content"
type TabPanel interface {
	// TabPanel is a Container.
	PanelView

	// TabBar returns the tab bar.
	TabBar() TabBar

	// TabBarPlacement returns the tab bar placement.
	TabBarPlacement() TabBarPlacement

	// SetTabBarPlacement sets tab bar placement.
	// Also sets the alignment of the tab bar according
	// to the tab bar placement.
	SetTabBarPlacement(tabBarPlacement TabBarPlacement)

	// TabBarFmt returns the cell formatter of the tab bar.
	TabBarFmt() CellFmt

	// Add adds a new tab (component) and an associated (content) component
	// to the tab panel.
	Add(tab, content Comp)

	// Add adds a new tab (string) and an associated (content) component
	// to the tab panel.
	// This is a shorthand for
	// 		Add(NewLabel(tab), content)
	AddString(tab string, content Comp)

	// Selected returns the selected tab idx.
	// Returns -1 if no tab is selected.
	Selected() int

	// SetSelected sets the selected tab idx.
	// If idx < 0, no tabs will be selected.
	// If idx > CompsCount(), this is a no-op.
	SetSelected(idx int)
}

// TabPanel implementation.
type tabPanelImpl struct {
	panelImpl // panel implementation: TabPanel is a Panel, but only PanelView's methods are exported.

	tabBarImpl      *tabBarImpl     // Tab bar implementation
	tabBarPlacement TabBarPlacement // Tab bar placement
	tabBarFmt       *cellFmtImpl    // Tab bar cell formatter

	selected int // The selected tab idx
}

// NewTabPanel creates a new TabPanel.
// Default tab bar placement is TB_PLACEMENT_TOP,
// default horizontal alignment is HA_DEFAULT,
// default vertical alignment is VA_DEFAULT.
func NewTabPanel() TabPanel {
	c := &tabPanelImpl{panelImpl: newPanelImpl(), tabBarImpl: newTabBarImpl(), tabBarFmt: newCellFmtImpl(), selected: -1}
	c.SetAttr("cellspacing", "0")
	c.SetAttr("cellpadding", "0")
	c.tabBarImpl.Style().AddClass("gwu-TabBar")
	c.tabBarImpl.setParent(c)
	c.SetTabBarPlacement(TB_PLACEMENT_TOP)
	c.tabBarFmt.SetAlign(HA_LEFT, VA_TOP)
	c.Style().AddClass("gwu-TabPanel")
	return c
}

func (c *tabPanelImpl) Remove(c2 Comp) bool {
	i := c.CompIdx(c2)
	if i < 0 {
		// Try the tab bar:
		i = c.tabBarImpl.CompIdx(c2)
		if i < 0 {
			return false
		}

		// It's a tab component
		return c.Remove(c.panelImpl.CompAt(i))
	}

	// It's a content component
	c.tabBarImpl.panelImpl.Remove(c.tabBarImpl.CompAt(i))
	c.panelImpl.Remove(c2)

	if i < c.selected {
		c.selected-- // Keep the same tab selected by decreasing its index by 1
	} else if i == c.selected { // Selected tab was removed...
		if i < c.CompsCount() {
			c.SetSelected(i) // There is next tab, select it
		} else if i > 0 { // Last was selected and removed...
			c.SetSelected(i - 1) // ...select the "new" last one
		} else {
			c.SetSelected(-1) // No tabs remained.
		}
	}

	return true
}

func (c *tabPanelImpl) ById(id ID) Comp {
	// panelImpl.ById() also checks our own id first
	c2 := c.panelImpl.ById(id)
	if c2 != nil {
		return c2
	}

	c2 = c.tabBarImpl.ById(id)
	if c2 != nil {
		return c2
	}

	return nil
}

func (c *tabPanelImpl) Clear() {
	c.tabBarImpl.Clear()
	c.panelImpl.Clear()

	c.SetSelected(-1)
}

func (c *tabPanelImpl) TabBar() TabBar {
	return c.tabBarImpl
}

func (c *tabPanelImpl) TabBarPlacement() TabBarPlacement {
	return c.tabBarPlacement
}

func (c *tabPanelImpl) SetTabBarPlacement(tabBarPlacement TabBarPlacement) {
	style := c.tabBarImpl.Style()

	// Remove old style class
	switch c.tabBarPlacement {
	case TB_PLACEMENT_TOP:
		style.RemoveClass("gwu-TabBar-Top")
	case TB_PLACEMENT_BOTTOM:
		style.RemoveClass("gwu-TabBar-Bottom")
	case TB_PLACEMENT_LEFT:
		style.RemoveClass("gwu-TabBar-Left")
	case TB_PLACEMENT_RIGHT:
		style.RemoveClass("gwu-TabBar-Right")
	}

	c.tabBarPlacement = tabBarPlacement

	switch tabBarPlacement {
	case TB_PLACEMENT_TOP:
		c.tabBarImpl.SetLayout(LAYOUT_HORIZONTAL)
		c.tabBarImpl.SetAlign(HA_LEFT, VA_BOTTOM)
		style.AddClass("gwu-TabBar-Top")
	case TB_PLACEMENT_BOTTOM:
		c.tabBarImpl.SetLayout(LAYOUT_HORIZONTAL)
		c.tabBarImpl.SetAlign(HA_LEFT, VA_TOP)
		style.AddClass("gwu-TabBar-Bottom")
	case TB_PLACEMENT_LEFT:
		c.tabBarImpl.SetLayout(LAYOUT_VERTICAL)
		c.tabBarImpl.SetAlign(HA_RIGHT, VA_TOP)
		style.AddClass("gwu-TabBar-Left")
	case TB_PLACEMENT_RIGHT:
		c.tabBarImpl.SetLayout(LAYOUT_VERTICAL)
		c.tabBarImpl.SetAlign(HA_LEFT, VA_TOP)
		style.AddClass("gwu-TabBar-Right")
	}
}

func (c *tabPanelImpl) TabBarFmt() CellFmt {
	return c.tabBarFmt
}

func (c *tabPanelImpl) Add(tab, content Comp) {
	c.tabBarImpl.Add(tab)
	c.panelImpl.Add(content)
	tab.Style().AddClass("gwu-TabBar-NotSelected")
	c.CellFmt(content).Style().AddClass("gwu-TabPanel-Content")

	if c.CompsCount() == 1 {
		c.SetSelected(0)
	}

	tab.AddEHandlerFunc(func(e Event) {
		c.SetSelected(c.CompIdx(content))
		e.MarkDirty(c)
	}, ETYPE_CLICK)
}

func (c *tabPanelImpl) AddString(tab string, content Comp) {
	c.Add(NewLabel(tab), content)
}

func (c *tabPanelImpl) Selected() int {
	return c.selected
}

func (c *tabPanelImpl) SetSelected(idx int) {
	if idx >= c.CompsCount() {
		return
	}

	if c.selected >= 0 {
		// Deselect current selected
		style := c.tabBarImpl.CompAt(c.selected).Style()
		style.RemoveClass("gwu-TabBar-Selected")
		style.AddClass("gwu-TabBar-NotSelected")
	}

	c.selected = idx

	if c.selected >= 0 {
		// Select new selected
		style := c.tabBarImpl.CompAt(idx).Style()
		style.RemoveClass("gwu-TabBar-NotSelected")
		style.AddClass("gwu-TabBar-Selected")
	}
}

func (c *tabPanelImpl) Render(w writer) {
	w.Write(_STR_TABLE_OP)
	c.renderAttrsAndStyle(w)
	c.renderEHandlers(w)
	w.Write(_STR_GT)

	switch c.tabBarPlacement {
	case TB_PLACEMENT_TOP:
		c.tabBarFmt.render("tr", w)
		w.Write(_STR_TD)
		c.tabBarImpl.Render(w)
		w.Writes(c.trTag())
		c.renderContent(w)
	case TB_PLACEMENT_BOTTOM:
		w.Writes(c.trTag())
		c.renderContent(w)
		c.tabBarFmt.render("tr", w)
		w.Write(_STR_TD)
		c.tabBarImpl.Render(w)
	case TB_PLACEMENT_LEFT:
		w.Writes(c.trTag())
		c.tabBarFmt.render("td", w)
		c.tabBarImpl.Render(w)
		c.renderContent(w)
	case TB_PLACEMENT_RIGHT:
		w.Writes(c.trTag())
		c.renderContent(w)
		c.tabBarFmt.render("td", w)
		c.tabBarImpl.Render(w)
	}

	w.Write(_STR_TABLE_CL)
}

// renderContent renders the selected content component.
func (c *tabPanelImpl) renderContent(w writer) {
	// Render only the selected content component
	if c.selected >= 0 {
		c2 := c.comps[c.selected]
		c.renderTd(c2, w)
		c2.Render(w)
	} else {
		w.Write(_STR_TD)
	}
}
