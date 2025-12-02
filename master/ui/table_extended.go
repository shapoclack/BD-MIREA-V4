package table

import (
	"BD_Mirea/internal"
	"context"
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ============ VIEW UI Functions ============

// UICreateView handles view creation in UI
func UICreateView(ctx context.Context, pool *pgxpool.Pool, window fyne.Window) {
	viewNameEntry := widget.NewEntry()
	viewNameEntry.SetPlaceHolder("my_view")

	selectQueryEntry := widget.NewMultiLineEntry()
	selectQueryEntry.SetPlaceHolder("SELECT id, name FROM products WHERE active = true")
	selectQueryEntry.SetMinRowsVisible(6)

	form := container.NewVBox(
		widget.NewLabel("View Name:"),
		viewNameEntry,
		widget.NewLabel("SELECT Query:"),
		selectQueryEntry,
	)

	dialog.ShowCustomConfirm("Create VIEW", "Create", "Cancel", form, func(ok bool) {
		if !ok {
			return
		}

		viewName := strings.TrimSpace(viewNameEntry.Text)
		selectQuery := strings.TrimSpace(selectQueryEntry.Text)

		if viewName == "" || selectQuery == "" {
			showError(window, "View name and SELECT query are required")
			return
		}

		err := internal.CreateView(ctx, pool, viewName, selectQuery)
		if err != nil {
			showError(window, fmt.Sprintf("Failed to create view: %v", err))
			return
		}

		showInfo(window, fmt.Sprintf("VIEW '%s' created successfully!", viewName))
	}, window)
}

// UICreateOrReplaceView handles view creation or replacement
func UICreateOrReplaceView(ctx context.Context, pool *pgxpool.Pool, window fyne.Window) {
	viewNameEntry := widget.NewEntry()
	viewNameEntry.SetPlaceHolder("my_view")

	selectQueryEntry := widget.NewMultiLineEntry()
	selectQueryEntry.SetPlaceHolder("SELECT id, name FROM products WHERE active = true")
	selectQueryEntry.SetMinRowsVisible(6)

	form := container.NewVBox(
		widget.NewLabel("View Name:"),
		viewNameEntry,
		widget.NewLabel("SELECT Query:"),
		selectQueryEntry,
	)

	dialog.ShowCustomConfirm("Create or Replace VIEW", "Create/Replace", "Cancel", form, func(ok bool) {
		if !ok {
			return
		}

		viewName := strings.TrimSpace(viewNameEntry.Text)
		selectQuery := strings.TrimSpace(selectQueryEntry.Text)

		if viewName == "" || selectQuery == "" {
			showError(window, "View name and SELECT query are required")
			return
		}

		err := internal.CreateOrReplaceView(ctx, pool, viewName, selectQuery)
		if err != nil {
			showError(window, fmt.Sprintf("Failed to create or replace view: %v", err))
			return
		}

		showInfo(window, fmt.Sprintf("VIEW '%s' created or updated successfully!", viewName))
	}, window)
}

// UIListViews displays all views
func UIListViews(ctx context.Context, pool *pgxpool.Pool, window fyne.Window) {
	views, err := internal.ListAllViews(ctx, pool)
	if err != nil {
		showError(window, fmt.Sprintf("Failed to list views: %v", err))
		return
	}

	var tableData [][]string
	tableData = append(tableData, []string{"View Name"})

	for _, v := range views {
		tableData = append(tableData, []string{v})
	}

	table, err := CreateTable(tableData)
	if err != nil {
		showError(window, fmt.Sprintf("Failed to create table: %v", err))
		return
	}

	viewsWindow := fyne.CurrentApp().NewWindow("All VIEWs")
	viewsWindow.SetTitle("All VIEWs")
	viewsWindow.SetContent(container.NewScroll(table))
	viewsWindow.Resize(fyne.NewSize(500, 400))
	viewsWindow.CenterOnScreen()
	viewsWindow.Show()
}

// UIGetViewDefinition retrieves and displays view definition
func UIGetViewDefinition(ctx context.Context, pool *pgxpool.Pool, window fyne.Window) {
	viewNameEntry := widget.NewEntry()
	viewNameEntry.SetPlaceHolder("view_name")

	form := widget.NewForm(
		widget.NewFormItem("View Name", viewNameEntry),
	)

	dialog.ShowCustomConfirm("View Definition", "Show", "Cancel", form, func(ok bool) {
		if !ok {
			return
		}

		viewName := strings.TrimSpace(viewNameEntry.Text)
		if viewName == "" {
			showError(window, "View name is required")
			return
		}

		definition, err := internal.GetViewDefinition(ctx, pool, viewName)
		if err != nil {
			showError(window, fmt.Sprintf("Failed to get view definition: %v", err))
			return
		}

		defLabel := widget.NewLabel(definition)
		defLabel.Wrapping = fyne.TextWrapWord

		infoWindow := fyne.CurrentApp().NewWindow("View Definition")
		infoWindow.SetTitle(fmt.Sprintf("Definition of %s", viewName))
		infoWindow.SetContent(container.NewScroll(container.NewVBox(
			widget.NewCard("VIEW Definition", "", defLabel),
		)))
		infoWindow.Resize(fyne.NewSize(700, 400))
		infoWindow.CenterOnScreen()
		infoWindow.Show()
	}, window)
}

// UIDropView handles view deletion
func UIDropView(ctx context.Context, pool *pgxpool.Pool, window fyne.Window) {
	viewNameEntry := widget.NewEntry()
	viewNameEntry.SetPlaceHolder("view_name")

	warningLabel := widget.NewLabel("⚠️ WARNING: This action cannot be undone! All dependent objects will be deleted.")
	warningLabel.Wrapping = fyne.TextWrapWord

	confirmEntry := widget.NewEntry()
	confirmEntry.SetPlaceHolder(fmt.Sprintf("Type 'DELETE' to confirm"))

	form := container.NewVBox(
		widget.NewForm(
			widget.NewFormItem("View Name", viewNameEntry),
		),
		widget.NewSeparator(),
		warningLabel,
		widget.NewForm(
			widget.NewFormItem("Confirmation", confirmEntry),
		),
	)

	dialog.ShowCustomConfirm("Drop VIEW", "Delete", "Cancel", form, func(ok bool) {
		if !ok {
			return
		}

		viewName := strings.TrimSpace(viewNameEntry.Text)
		confirmation := strings.TrimSpace(confirmEntry.Text)

		if viewName == "" {
			showError(window, "View name is required")
			return
		}

		if confirmation != "DELETE" {
			showError(window, "Please type 'DELETE' to confirm deletion")
			return
		}

		err := internal.DropView(ctx, pool, viewName)
		if err != nil {
			showError(window, fmt.Sprintf("Failed to drop view: %v", err))
			return
		}

		showInfo(window, fmt.Sprintf("VIEW '%s' dropped successfully!", viewName))
	}, window)
}

// ============ MATERIALIZED VIEW UI Functions ============

// UICreateMaterializedView creates a materialized view
func UICreateMaterializedView(ctx context.Context, pool *pgxpool.Pool, window fyne.Window) {
	mvNameEntry := widget.NewEntry()
	mvNameEntry.SetPlaceHolder("my_materialized_view")

	selectQueryEntry := widget.NewMultiLineEntry()
	selectQueryEntry.SetPlaceHolder("SELECT id, name, COUNT(*) as cnt FROM products GROUP BY id, name")
	selectQueryEntry.SetMinRowsVisible(6)

	form := container.NewVBox(
		widget.NewLabel("Materialized View Name:"),
		mvNameEntry,
		widget.NewLabel("SELECT Query:"),
		selectQueryEntry,
	)

	dialog.ShowCustomConfirm("Create MATERIALIZED VIEW", "Create", "Cancel", form, func(ok bool) {
		if !ok {
			return
		}

		mvName := strings.TrimSpace(mvNameEntry.Text)
		selectQuery := strings.TrimSpace(selectQueryEntry.Text)

		if mvName == "" || selectQuery == "" {
			showError(window, "MV name and SELECT query are required")
			return
		}

		err := internal.CreateMaterializedView(ctx, pool, mvName, selectQuery)
		if err != nil {
			showError(window, fmt.Sprintf("Failed to create materialized view: %v", err))
			return
		}

		showInfo(window, fmt.Sprintf("MATERIALIZED VIEW '%s' created successfully!", mvName))
	}, window)
}

// UIRefreshMaterializedView refreshes a materialized view
func UIRefreshMaterializedView(ctx context.Context, pool *pgxpool.Pool, window fyne.Window) {
	mvNameEntry := widget.NewEntry()
	mvNameEntry.SetPlaceHolder("materialized_view_name")

	concurrentlyCheck := widget.NewCheck("Refresh CONCURRENTLY (if unique index exists)", nil)

	form := container.NewVBox(
		widget.NewForm(
			widget.NewFormItem("Materialized View Name", mvNameEntry),
		),
		widget.NewSeparator(),
		concurrentlyCheck,
	)

	dialog.ShowCustomConfirm("Refresh MATERIALIZED VIEW", "Refresh", "Cancel", form, func(ok bool) {
		if !ok {
			return
		}

		mvName := strings.TrimSpace(mvNameEntry.Text)
		if mvName == "" {
			showError(window, "MV name is required")
			return
		}

		err := internal.RefreshMaterializedView(ctx, pool, mvName, concurrentlyCheck.Checked)
		if err != nil {
			showError(window, fmt.Sprintf("Failed to refresh materialized view: %v", err))
			return
		}

		showInfo(window, fmt.Sprintf("MATERIALIZED VIEW '%s' refreshed successfully!", mvName))
	}, window)
}

// UIListMaterializedViews displays all materialized views
func UIListMaterializedViews(ctx context.Context, pool *pgxpool.Pool, window fyne.Window) {
	mvs, err := internal.ListAllMaterializedViews(ctx, pool)
	if err != nil {
		showError(window, fmt.Sprintf("Failed to list materialized views: %v", err))
		return
	}

	var tableData [][]string
	tableData = append(tableData, []string{"Materialized View Name"})

	for _, mv := range mvs {
		tableData = append(tableData, []string{mv})
	}

	table, err := CreateTable(tableData)
	if err != nil {
		showError(window, fmt.Sprintf("Failed to create table: %v", err))
		return
	}

	mvsWindow := fyne.CurrentApp().NewWindow("All MATERIALIZED VIEWs")
	mvsWindow.SetTitle("All MATERIALIZED VIEWs")
	mvsWindow.SetContent(container.NewScroll(table))
	mvsWindow.Resize(fyne.NewSize(500, 400))
	mvsWindow.CenterOnScreen()
	mvsWindow.Show()
}

// UIDropMaterializedView drops a materialized view
func UIDropMaterializedView(ctx context.Context, pool *pgxpool.Pool, window fyne.Window) {
	mvNameEntry := widget.NewEntry()
	mvNameEntry.SetPlaceHolder("materialized_view_name")

	warningLabel := widget.NewLabel("⚠️ WARNING: This action cannot be undone! All cached data will be deleted.")
	warningLabel.Wrapping = fyne.TextWrapWord

	confirmEntry := widget.NewEntry()
	confirmEntry.SetPlaceHolder("Type 'DELETE' to confirm")

	form := container.NewVBox(
		widget.NewForm(
			widget.NewFormItem("Materialized View Name", mvNameEntry),
		),
		widget.NewSeparator(),
		warningLabel,
		widget.NewForm(
			widget.NewFormItem("Confirmation", confirmEntry),
		),
	)

	dialog.ShowCustomConfirm("Drop MATERIALIZED VIEW", "Delete", "Cancel", form, func(ok bool) {
		if !ok {
			return
		}

		mvName := strings.TrimSpace(mvNameEntry.Text)
		confirmation := strings.TrimSpace(confirmEntry.Text)

		if mvName == "" {
			showError(window, "MV name is required")
			return
		}

		if confirmation != "DELETE" {
			showError(window, "Please type 'DELETE' to confirm deletion")
			return
		}

		err := internal.DropMaterializedView(ctx, pool, mvName)
		if err != nil {
			showError(window, fmt.Sprintf("Failed to drop materialized view: %v", err))
			return
		}

		showInfo(window, fmt.Sprintf("MATERIALIZED VIEW '%s' dropped successfully!", mvName))
	}, window)
}

// ============ ROLLUP/CUBE/GROUPING SETS UI Functions ============

// UIRollupQuery handles ROLLUP aggregation
func UIRollupQuery(ctx context.Context, pool *pgxpool.Pool, window fyne.Window) {
	tableEntry := widget.NewEntry()
	tableEntry.SetPlaceHolder("products")

	columnsEntry := widget.NewMultiLineEntry()
	columnsEntry.SetPlaceHolder("category, year, month")
	columnsEntry.SetMinRowsVisible(3)

	aggregateFunc := widget.NewSelect([]string{"SUM", "COUNT", "AVG", "MIN", "MAX"}, nil)
	aggregateFunc.SetSelected("SUM")

	aggregateColumnEntry := widget.NewEntry()
	aggregateColumnEntry.SetPlaceHolder("price")

	form := container.NewVBox(
		widget.NewForm(
			widget.NewFormItem("Table", tableEntry),
			widget.NewFormItem("Group By Columns", columnsEntry),
			widget.NewFormItem("Aggregate Function", aggregateFunc),
			widget.NewFormItem("Aggregate Column", aggregateColumnEntry),
		),
	)

	dialog.ShowCustomConfirm("ROLLUP Aggregation", "Execute", "Cancel", form, func(ok bool) {
		if !ok {
			return
		}

		table := strings.TrimSpace(tableEntry.Text)
		if table == "" {
			showError(window, "Table name is required")
			return
		}

		columnsText := strings.TrimSpace(columnsEntry.Text)
		if columnsText == "" {
			showError(window, "At least one grouping column is required")
			return
		}

		columns := strings.Split(columnsText, ",")
		for i := range columns {
			columns[i] = strings.TrimSpace(columns[i])
		}

		aggFunc := aggregateFunc.Selected
		aggColumn := strings.TrimSpace(aggregateColumnEntry.Text)

		results, err := internal.ExecuteRollupQuery(ctx, pool, table, columns, aggFunc, aggColumn)
		if err != nil {
			showError(window, fmt.Sprintf("Failed to execute ROLLUP query: %v", err))
			return
		}

		resultTable, err := CreateTable(results)
		if err != nil {
			showError(window, fmt.Sprintf("Failed to display results: %v", err))
			return
		}

		resultWindow := fyne.CurrentApp().NewWindow("ROLLUP Results")
		resultWindow.SetTitle("ROLLUP Results")
		resultWindow.SetContent(container.NewScroll(resultTable))
		resultWindow.Resize(fyne.NewSize(900, 600))
		resultWindow.CenterOnScreen()
		resultWindow.Show()
	}, window)
}

// UICubeQuery handles CUBE aggregation
func UICubeQuery(ctx context.Context, pool *pgxpool.Pool, window fyne.Window) {
	tableEntry := widget.NewEntry()
	tableEntry.SetPlaceHolder("products")

	columnsEntry := widget.NewMultiLineEntry()
	columnsEntry.SetPlaceHolder("category, year, month")
	columnsEntry.SetMinRowsVisible(3)

	aggregateFunc := widget.NewSelect([]string{"SUM", "COUNT", "AVG", "MIN", "MAX"}, nil)
	aggregateFunc.SetSelected("SUM")

	aggregateColumnEntry := widget.NewEntry()
	aggregateColumnEntry.SetPlaceHolder("price")

	form := container.NewVBox(
		widget.NewForm(
			widget.NewFormItem("Table", tableEntry),
			widget.NewFormItem("Group By Columns", columnsEntry),
			widget.NewFormItem("Aggregate Function", aggregateFunc),
			widget.NewFormItem("Aggregate Column", aggregateColumnEntry),
		),
	)

	dialog.ShowCustomConfirm("CUBE Aggregation", "Execute", "Cancel", form, func(ok bool) {
		if !ok {
			return
		}

		table := strings.TrimSpace(tableEntry.Text)
		if table == "" {
			showError(window, "Table name is required")
			return
		}

		columnsText := strings.TrimSpace(columnsEntry.Text)
		if columnsText == "" {
			showError(window, "At least one grouping column is required")
			return
		}

		columns := strings.Split(columnsText, ",")
		for i := range columns {
			columns[i] = strings.TrimSpace(columns[i])
		}

		aggFunc := aggregateFunc.Selected
		aggColumn := strings.TrimSpace(aggregateColumnEntry.Text)

		results, err := internal.ExecuteCubeQuery(ctx, pool, table, columns, aggFunc, aggColumn)
		if err != nil {
			showError(window, fmt.Sprintf("Failed to execute CUBE query: %v", err))
			return
		}

		resultTable, err := CreateTable(results)
		if err != nil {
			showError(window, fmt.Sprintf("Failed to display results: %v", err))
			return
		}

		resultWindow := fyne.CurrentApp().NewWindow("Window")
		resultWindow.SetTitle("CUBE Results")
		resultWindow.SetContent(container.NewScroll(resultTable))
		resultWindow.Resize(fyne.NewSize(900, 600))
		resultWindow.CenterOnScreen()
		resultWindow.Show()
	}, window)
}

// ============ CTE (WITH) UI Functions ============

// UICTEBuilder handles Common Table Expression creation
func UICTEBuilder(ctx context.Context, pool *pgxpool.Pool, window fyne.Window) {
	cteNameEntry := widget.NewEntry()
	cteNameEntry.SetPlaceHolder("cte_name")

	cteQueryEntry := widget.NewMultiLineEntry()
	cteQueryEntry.SetPlaceHolder("SELECT id, name FROM products WHERE price > 100")
	cteQueryEntry.SetMinRowsVisible(4)

	mainQueryEntry := widget.NewMultiLineEntry()
	mainQueryEntry.SetPlaceHolder("SELECT * FROM cte_name WHERE id > 5")
	mainQueryEntry.SetMinRowsVisible(4)

	form := container.NewVBox(
		widget.NewLabel("CTE Definition:"),
		widget.NewForm(
			widget.NewFormItem("CTE Name", cteNameEntry),
		),
		widget.NewLabel("CTE Query (SELECT):"),
		cteQueryEntry,
		widget.NewLabel("Main Query (using CTE):"),
		mainQueryEntry,
	)

	dialog.ShowCustomConfirm("WITH (CTE) Query", "Execute", "Cancel", form, func(ok bool) {
		if !ok {
			return
		}

		cteName := strings.TrimSpace(cteNameEntry.Text)
		cteQuery := strings.TrimSpace(cteQueryEntry.Text)
		mainQuery := strings.TrimSpace(mainQueryEntry.Text)

		if cteName == "" || cteQuery == "" || mainQuery == "" {
			showError(window, "All fields are required")
			return
		}

		// Extract table name from main query for QueryBuilder
		tableFromQuery := "products" // default
		parts := strings.Fields(mainQuery)
		for i, part := range parts {
			if strings.ToUpper(part) == "FROM" && i+1 < len(parts) {
				tableFromQuery = parts[i+1]
				break
			}
		}

		cteDefinitions := []internal.CTEDefinition{
			{
				Name:    cteName,
				Query:   cteQuery,
				Columns: []string{},
			},
		}

		mainQB := internal.NewQueryBuilder(tableFromQuery)
		mainQB.Where(mainQuery)

		results, err := internal.ExecuteCTEQuery(ctx, pool, cteDefinitions, mainQB)
		if err != nil {
			showError(window, fmt.Sprintf("Failed to execute CTE query: %v", err))
			return
		}

		resultTable, err := CreateTable(results)
		if err != nil {
			showError(window, fmt.Sprintf("Failed to display results: %v", err))
			return
		}

		resultWindow := fyne.CurrentApp().NewWindow("Window")
		resultWindow.SetTitle("CTE (WITH) Results")
		resultWindow.SetContent(container.NewScroll(resultTable))
		resultWindow.Resize(fyne.NewSize(900, 600))
		resultWindow.CenterOnScreen()
		resultWindow.Show()
	}, window)
}
