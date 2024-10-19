import React, { Component } from 'react';

import CodeEditor from '@uiw/react-textarea-code-editor';
import rehypePrism from 'rehype-prism-plus';

import Spinner from "react-bootstrap/Spinner";
import Nav from 'react-bootstrap/Nav';
import FloatingLabel from "react-bootstrap/FloatingLabel";

import Input from "./input.component";

import { control } from "react-validation";

import { required } from "./validations";

import IntegrationService from "../../services/integration.service";

import StorageRounded from '@mui/icons-material/StorageRounded';
import TableChartRounded from '@mui/icons-material/TableChartRounded';
import AddRounded from '@mui/icons-material/AddRounded';
import DeleteOutlineRounded from '@mui/icons-material/DeleteOutlineRounded';
import ClearRounded from '@mui/icons-material/ClearRounded';

import {
    Dropdown,
    Button,
} from 'react-bootstrap'

export const orderDirectionToHuman = {
  nat:  'natural',
  asc:  'ascending',
  desc: 'descending',
}

export const orderDirectionToSQL = {
  asc:  'ASC',
  desc: 'DESC',
}

export const equationToHuman = {
  neq:  'is not equal to',
  let:  'is less than or equal to',
  het:  'is greater than or equal to',
  eq:   'is equal to',
  lt:   'is less than',
  ht:   'is greater than',
  ct:   'contains',
  isn:  'is null',
  isnn: 'is not null',
  ise:  'is empty',
  isne: 'is not empty',
  ltsa: 'is less than N seconds ago',
  ltma: 'is less than N minutes ago',
  ltha: 'is less than N hours ago',
  ltda: 'is less than N days ago',
  mtsa: 'is more than N seconds ago',
  mtma: 'is more than N minutes ago',
  mtha: 'is more than N hours ago',
  mtda: 'is more than N days ago',
  curd: 'is today',
  curm: 'is in current month',
  cury: 'is in current year',
};

export const equationToHumanRegex = {
  ltsa: "is less than $1 seconds ago",
  ltma: "is less than $1 minutes ago",
  ltha: "is less than $1 hours ago",
  ltda: "is less than $1 days ago",
  mtsa: "is more than $1 seconds ago",
  mtma: "is more than $1 minutes ago",
  mtha: "is more than $1 hours ago",
  mtda: "is more than $1 days ago",
  curd: '$1 is today',
  curm: '$1 is in current month',
  cury: '$1 is in current year',
};

export const equationToSQL = {
  ltsa: "> DATE_SUB\\(NOW\\(\\), INTERVAL (\\d*) SECOND\\)",
  ltma: "> DATE_SUB\\(NOW\\(\\), INTERVAL (\\d*) MINUTE\\)",
  ltha: "> DATE_SUB\\(NOW\\(\\), INTERVAL (\\d*) HOUR\\)",
  ltda: "> DATE_SUB\\(NOW\\(\\), INTERVAL (\\d*) DAY\\)",
  mtsa: "< DATE_SUB\\(NOW\\(\\), INTERVAL (\\d*) SECOND\\)",
  mtma: "< DATE_SUB\\(NOW\\(\\), INTERVAL (\\d*) MINUTE\\)",
  mtha: "< DATE_SUB\\(NOW\\(\\), INTERVAL (\\d*) HOUR\\)",
  mtda: "< DATE_SUB\\(NOW\\(\\), INTERVAL (\\d*) DAY\\)",
  curd: 'YEAR\\(([^\\)]*)\\) = YEAR\\(CURRENT_DATE\\(\\)\\)[ \n]+AND MONTH\\([^\\)]*\\) = MONTH\\(CURRENT_DATE\\(\\)\\)[ \n]+AND DAY\\([^\\)]*\\) = DAY\\(CURRENT_DATE\\(\\)\\)',
  curm: 'YEAR\\(([^\\)]*)\\) = YEAR\\(CURRENT_DATE\\(\\)\\)[ \n]+AND MONTH\\([^\\)]*\\) = MONTH\\(CURRENT_DATE\\(\\)\\)',
  cury: 'YEAR\\(([^\\)]*)\\) = YEAR\\(CURRENT_DATE\\(\\)\\)',
  isne: '!= ""',
  ise:  '= ""',
  neq:  '!=',
  let:  '<=',
  het:  '>=',
  eq:   '=',
  lt:   '<',
  ht:   '>',
  ct:   'LIKE',
  isn:  'IS NULL',
  isnn: 'IS NOT NULL',
};

export const filterToSQL = {
  neq:  "%%%key%%% != %%%value%%%",
  let:  "%%%key%%% <= %%%value%%%",
  het:  "%%%key%%% >= %%%value%%%",
  eq:   "%%%key%%% = %%%value%%%",
  lt:   "%%%key%%% < %%%value%%%",
  ht:   "%%%key%%% > %%%value%%%",
  ct:   '%%%key%%% LIKE "%%%%value%%%%"',
  isn:  "%%%key%%% IS NULL",
  isnn: "%%%key%%% IS NOT NULL",
  isne: '%%%key%%% != ""',
  ise:  '%%%key%%% = ""',
  ltsa: '%%%key%%% > DATE_SUB(NOW(), INTERVAL %%%value%%% SECOND)',
  ltma: '%%%key%%% > DATE_SUB(NOW(), INTERVAL %%%value%%% MINUTE)',
  ltha: '%%%key%%% > DATE_SUB(NOW(), INTERVAL %%%value%%% HOUR)',
  ltda: '%%%key%%% > DATE_SUB(NOW(), INTERVAL %%%value%%% DAY)',
  mtsa: '%%%key%%% < DATE_SUB(NOW(), INTERVAL %%%value%%% SECOND)',
  mtma: '%%%key%%% < DATE_SUB(NOW(), INTERVAL %%%value%%% MINUTE)',
  mtha: '%%%key%%% < DATE_SUB(NOW(), INTERVAL %%%value%%% HOUR)',
  mtda: '%%%key%%% < DATE_SUB(NOW(), INTERVAL %%%value%%% DAY)',
  curd: 'YEAR(%%%key%%%) = YEAR(CURRENT_DATE()) AND MONTH(%%%key%%%) = MONTH(CURRENT_DATE()) AND DAY(%%%key%%%) = DAY(CURRENT_DATE())',
  curm: 'YEAR(%%%key%%%) = YEAR(CURRENT_DATE()) AND MONTH(%%%key%%%) = MONTH(CURRENT_DATE())',
  cury: 'YEAR(%%%key%%%) = YEAR(CURRENT_DATE())',
}

export const equationNoValue = {
  isn:  '',
  isnn: '',
  curd: '',
  curm: '',
  cury: '',
  ise:  '',
  isne: '',
}

export const equationContains = {
  ct: '',
}

class QueryUI extends Component {
    constructor(props) {
      super(props);

      this.onChangeQuery = this.onChangeQuery.bind(this);
      this.onChangeDatabase = this.onChangeDatabase.bind(this);
      this.onChangeTable = this.onChangeTable.bind(this);
      this.onChangeFilterColumn = this.onChangeFilterColumn.bind(this);
      this.onChangeGroupByColumn = this.onChangeGroupByColumn.bind(this);
      this.onChangeOrderByColumn = this.onChangeOrderByColumn.bind(this);
      this.onChangeOrderDirection = this.onChangeOrderDirection.bind(this);
      this.onChangeFilterOperation = this.onChangeFilterOperation.bind(this);
      this.onChangeFilterValue = this.onChangeFilterValue.bind(this);
      this.toggleSQLForm = this.toggleSQLForm.bind(this);
      this.toggleFilterForm = this.toggleFilterForm.bind(this);
      this.editorToSQL = this.editorToSQL.bind(this);
      this.SQLToEditor = this.SQLToEditor.bind(this);
      this.isSQLEditorCompatible = this.isSQLEditorCompatible.bind(this);
      this.onChangeQueryFromEditor = this.onChangeQueryFromEditor.bind(this);
      this.filterToHumanReadable = this.filterToHumanReadable.bind(this);
      this.orderByToHumanReadable = this.orderByToHumanReadable.bind(this);
      this.isFilterEditable = this.isFilterEditable.bind(this);
      this.removeFilter = this.removeFilter.bind(this);
      this.addFilter = this.addFilter.bind(this);
      this.toggleGroupByForm = this.toggleGroupByForm.bind(this);
      this.removeGroupBy = this.removeGroupBy.bind(this);
      this.addGroupBy = this.addGroupBy.bind(this);
      this.toggleOrderByForm = this.toggleOrderByForm.bind(this);
      this.removeOrderBy = this.removeOrderBy.bind(this);
      this.addOrderBy = this.addOrderBy.bind(this);

      this.state = {
        query: this.props.query,
        databases: null,
        database: null,
        tables: null,
        table: null,
        filterColumn: null,
        filterOperation: equationToHuman.neq,
        filterValue: "",
        filterOperationSql: equationToSQL.neq,
        filterOperationKey: "neq",
        filters: [],
        humanFilters: [],
        groupBys: [],
        groupByColumn: null,
        orderBys: [],
        humanOrderBys: [],
        orderByColumn: null,
        orderDirection: orderDirectionToHuman.nat,
        orderDirectionSql: "",
        orderDirectionKey: "nat",
        columns: null,
        integrationUuid: null,
        isFilterFormOpen: false,
        isGroupByFormOpen: false,
        isOrderByFormOpen: false,
        showSQLForm: localStorage.getItem('showSQLForm') === "true",
      };
    }

    componentDidMount = async() => {
        await this.promisedSetState({
            query: this.props.query,
            databases: null,
            database: null,
            tables: null,
            table: null,
            filterColumn: null,
            filterOperation: equationToHuman.neq,
            filterValue: "",
            filterOperationSql: equationToSQL.neq,
            filterOperationKey: "neq",
            filters: [],
            humanFilters: [],
            groupBys: [],
            groupByColumn: null,
            orderBys: [],
            humanOrderBys: [],
            orderByColumn: null,
            orderDirection: orderDirectionToHuman.nat,
            orderDirectionSql: "",
            orderDirectionKey: "nat",
            columns: null,
            integrationUuid: this.props.integrationUuid,
            showSQLForm: localStorage.getItem('showSQLForm') === "true",
        });

        if (this.state.showSQLForm === false) {
            if (this.isSQLEditorCompatible(this.props.query)) {
                await this.handleGetDatabases(this.state.integrationUuid);
                await this.SQLToEditor(this.props.query);
            }
        }
    };

    promisedSetState = (newState) => new Promise(resolve => this.setState(newState, resolve));

    UNSAFE_componentWillReceiveProps = async(nextProps, nextContext) => {
        if (this.state.integrationUuid !== this.props.integrationUuid) {
            await this.promisedSetState({
                query: this.props.query,
                databases: null,
                tables: null,
                columns: null,
                database: null,
                filterColumn: null,
                filterOperation: equationToHuman.neq,
                filterValue: "",
                filterOperationSql: equationToSQL.neq,
                filterOperationKey: "neq",
                table: null,
                filters: [],
                humanFilters: [],
                groupBys: [],
                groupByColumn: null,
                orderBys: [],
                humanOrderBys: [],
                orderByColumn: null,
                orderDirection: orderDirectionToHuman.nat,
                orderDirectionSql: "",
                orderDirectionKey: "nat",
            });

            if (this.state.showSQLForm === false) {
                if (this.isSQLEditorCompatible(this.props.query)) {
                    await this.handleGetDatabases(this.props.integrationUuid);
                    await this.SQLToEditor(this.props.query);
                }
            }

            await this.promisedSetState({
                integrationUuid: this.props.integrationUuid,
            });
        }
    }

    filterToHumanReadable(filter) {
        let regex = new RegExp("");
        let replacement = "";

        for (const [key, value] of Object.entries(equationToSQL)) {
            regex = new RegExp(value, "gi");
            if (key in equationToHumanRegex) {
                replacement = equationToHumanRegex[key];
            } else {
                replacement = equationToHuman[key];
            }

            filter = filter.replace(regex, " " + replacement + " ");
            filter = filter.replace("  ", " ");
        }

        return filter;
    }

    orderByToHumanReadable(order) {
        let regex = new RegExp("");
        let replacement = "";

        for (const [key, value] of Object.entries(orderDirectionToSQL)) {
            regex = new RegExp(value, "gi");
            replacement = orderDirectionToHuman[key];

            order = order.replace(regex, " " + replacement + " ");
            order = order.replace("  ", " ");
        }

        return order;
    }

    removeFilter = async(index) => {
        let filters = this.state.filters;
        let humanFilters = this.state.humanFilters;

        filters.splice(index, 1);
        humanFilters.splice(index, 1);

        await this.promisedSetState({filters, humanFilters});

        await this.onChangeQueryFromEditor();
    }

    addFilter = async() => {
        if (
            this.state.filterColumn !== null
        ) {
            let filters = this.state.filters;
            let humanFilters = this.state.humanFilters;

            let key = this.state.filterColumn;
            let value = this.state.filterValue;

            if (isNaN(this.state.filterValue) || value === "") {
                value = '"'  + value + '"';
            }
            
            let newFilter = filterToSQL[this.state.filterOperationKey].replaceAll("%%%key%%%", key);
            newFilter = newFilter.replaceAll('%%%value%%%', value)

            filters.push(newFilter);
            humanFilters.push(this.filterToHumanReadable(newFilter));

            await this.promisedSetState({filters, humanFilters});

            await this.onChangeQueryFromEditor();
        }

        this.toggleFilterForm();
    }

    removeGroupBy = async(index) => {
        let groupBys = this.state.groupBys;

        groupBys.splice(index, 1);

        await this.promisedSetState({groupBys});

        await this.onChangeQueryFromEditor();
    }

    addGroupBy = async() => {
        if (
            this.state.groupByColumn !== null
        ) {
            let groupBys = this.state.groupBys;

            let value = this.state.groupByColumn;
            
            groupBys.push(value);

            await this.promisedSetState({groupBys});

            await this.onChangeQueryFromEditor();
        }

        this.toggleGroupByForm();
    }

    removeOrderBy = async(index) => {
        let orderBys = this.state.orderBys;
        let humanOrderBys = this.state.humanOrderBys;

        orderBys.splice(index, 1);
        humanOrderBys.splice(index, 1);

        await this.promisedSetState({orderBys, humanOrderBys});

        await this.onChangeQueryFromEditor();
    }

    addOrderBy = async() => {
        if (
            this.state.orderByColumn !== null
        ) {
            let orderBys = this.state.orderBys;
            let humanOrderBys = this.state.humanOrderBys;

            let newOrderBy = this.state.orderByColumn;

            if (this.state.orderDirectionKey !== "nat") {
                newOrderBy = newOrderBy + " " + this.state.orderDirectionSql;
            }
            
            orderBys.push(newOrderBy);
            humanOrderBys.push(this.orderByToHumanReadable(newOrderBy));

            await this.promisedSetState({orderBys, humanOrderBys});

            await this.onChangeQueryFromEditor();
        }

        this.toggleOrderByForm();
    }

    isFilterEditable(filter) {
        let humanisedFilter = this.filterToHumanReadable(filter);
        
        return humanisedFilter !== filter;
    }

    editorToSQL() {
        let sql = "SELECT *\nFROM " + this.state.database + "." + this.state.table;
        
        if (this.state.filters.length > 0) {
            sql = sql + "\nWHERE";

            for (let i = 0; i < this.state.filters.length; i++) {
                if (i === 0) {
                    sql = sql + " " + this.state.filters[i].trim();
                } else {
                    sql = sql + "  AND " + this.state.filters[i].trim();
                }

                if (i !== this.state.filters.length - 1) {
                    sql = sql + "\n";
                }
            }
        }

        if (this.state.groupBys.length > 0) {
            sql = sql + "\nGROUP BY ";

            for (let i = 0; i < this.state.groupBys.length; i++) {
                if (i !== 0) {
                    sql = sql + ", ";
                }

                sql = sql + this.state.groupBys[i].trim();
            }
        }

        if (this.state.orderBys.length > 0) {
            sql = sql + "\nORDER BY ";

            for (let i = 0; i < this.state.orderBys.length; i++) {
                if (i !== 0) {
                    sql = sql + ", ";
                }

                sql = sql + this.state.orderBys[i].trim();
            }
        }

        sql = sql + "\n;";

        return sql;
    }

    SQLToEditor = async(sql) => {
        let database = null;
        let table = null;
        let filters = [];
        let humanFilters = [];
        let groupBys = [];
        let orderBys = [];
        let humanOrderBys = [];

        const inclusion = new RegExp(/SELECT[ |\n]+\*[ |\n]+FROM[ |\n]+([^.; \n]*)(\.)?([^.; \n]*)/i);
        let matches = inclusion.exec(sql);

        if (matches !== null) {
            if (matches[1] !== "") {
                if (matches[3] === "" ) {
                    table = matches[1]
                } else {
                    database = matches[1]
                }
            }
            if (matches[3] !== "") {
                table = matches[3]
            }
        }

        const filterSMatchingOne = new RegExp(/WHERE[ |\n]+([^;]*)/i);
        let filterStringOne = filterSMatchingOne.exec(sql);
        let filterString = "";
        if (filterStringOne !== null && filterStringOne[1] !== "") {
            const filterSMatchingTwo = new RegExp(/^([\S\s]*)GROUP BY/i);
            let filterStringTwo = filterSMatchingTwo.exec(filterStringOne[1]);

            if (filterStringTwo !== null && filterStringTwo[1] !== "") {
                filterString = filterStringTwo[1];
            } else {
                filterString = filterStringOne[1];
            }

            const filterSMatchingThree = new RegExp(/^([\S\s]*)ORDER BY/i);
            let filterStringThree = filterSMatchingThree.exec(filterString);
            if (filterStringThree !== null && filterStringThree[1] !== "") {
                filterString = filterStringThree[1];
            }

            filters = filterString.split(new RegExp("and", 'i'));

            let rawHumanFilters = this.filterToHumanReadable(filterString);
            humanFilters = rawHumanFilters.split(new RegExp("and", 'i'));
        }

        const groupBysSMatchingOne = new RegExp(/GROUP BY[ |\n]+([^;]*)/i);
        let groupBysStringOne = groupBysSMatchingOne.exec(sql);
        let groupBysString = "";
        if (groupBysStringOne !== null && groupBysStringOne[1] !== "") {
            const groupBysSMatchingTwo = new RegExp(/^([\S\s]*)ORDER BY/i);
            let groupBysStringTwo = groupBysSMatchingTwo.exec(groupBysStringOne[1]);

            if (groupBysStringTwo !== null && groupBysStringTwo[1] !== "") {
                groupBysString = groupBysStringTwo[1];
            } else {
                groupBysString = groupBysStringOne[1];
            }

            groupBys = groupBysString.split(new RegExp(",", 'i'));
        }

        const orderBysSMatching = new RegExp(/ORDER BY[ |\n]+([^;]*)/i);
        let orderBysString = orderBysSMatching.exec(sql);
        if (orderBysString !== null && orderBysString[1] !== "") {
            orderBys = orderBysString[1].split(new RegExp(",", 'i'));

            let rawHumanOrderBys = this.orderByToHumanReadable(orderBysString[1]);
            humanOrderBys = rawHumanOrderBys.split(new RegExp(",", 'i'));
        }

        await this.promisedSetState({database, table, filters, humanFilters, groupBys, orderBys, humanOrderBys});
    }

    toggleSQLForm = async(show) => {
        localStorage.setItem('showSQLForm', show);
        await this.promisedSetState({
            showSQLForm: show,
        });

        if (this.isSQLEditorCompatible(this.state.query)) {
            if (show === false) {
                this.handleGetDatabases(this.props.integrationUuid);
                this.SQLToEditor(this.state.query);
            } else {
                this.onChangeQueryFromEditor();
            }
        }
    }

    toggleFilterForm = async() => {
        await this.promisedSetState({
            isFilterFormOpen: !this.state.isFilterFormOpen,
            filterOperation: equationToHuman.neq,
            filterValue: "",
            filterOperationSql: equationToSQL.neq,
            filterOperationKey: "neq",
        });
    }

    toggleGroupByForm = async() => {
        await this.promisedSetState({
            isGroupByFormOpen: !this.state.isGroupByFormOpen,
        });
    }

    toggleOrderByForm = async() => {
        await this.promisedSetState({
            isOrderByFormOpen: !this.state.isOrderByFormOpen,
            orderDirection: orderDirectionToHuman.nat,
            orderDirectionSql: "",
            orderDirectionKey: "nat",
        });
    }

    handleGetDatabases = async(integrationUuid) => {
        let databases = this.state.databases;

        if (
            databases === null
            || databases.length === 0
        ) {
            databases = IntegrationService.getDatabases(integrationUuid);

            Promise.resolve(databases)
                .then(async(databases) => {
                    if (databases.data && databases.data.items) {
                        let dbs = databases.data.items;
                        let index = dbs.indexOf("information_schema");
                        if (index !== -1) {
                          dbs.splice(index, 1);
                        }
                        await this.promisedSetState({databases: dbs});
                        if (dbs.length > 0) {
                            this.onChangeDatabase(
                                this.state.database === null
                                    ? dbs[0]
                                    : this.state.database
                            );
                        }
                    } else {
                        await this.promisedSetState({databases: []});
                    }
                })
                .catch(async() => {
                    await this.promisedSetState({databases: []});
                });
        }
    };

    handleGetTables = async(integrationUuid, db) => {
        let tables = this.state.tables;

        if (
            tables === null
            || tables.length === 0
        ) {
            tables = IntegrationService.getTables(integrationUuid, db);

            Promise.resolve(tables)
                .then(async(tables) => {
                    if (tables.data && tables.data.items) {
                        await this.promisedSetState({tables: tables.data.items});
                        if (tables.data.items.length > 0) {                
                            this.onChangeTable(
                                this.state.table === null
                                    ? tables.data.items[0]
                                    : this.state.table
                            );
                        }
                    } else {
                        await this.promisedSetState({tables: []});
                    }
                })
                .catch(async() => {
                    await this.promisedSetState({tables: []});
                });
        }
    };

    handleDescribeTable = async(integrationUuid, db, table) => {
        let columns = this.state.columns;

        columns = IntegrationService.describeTable(integrationUuid, db, table);

        Promise.resolve(columns)
            .then(async(columns) => {
                if (columns.data && columns.data.items) {
                    await this.promisedSetState({columns: columns.data.items});
                    if (this.state.filterColumn === null) {
                        await this.promisedSetState({filterColumn: columns.data.items[0]}); 
                    }
                    if (this.state.groupByColumn === null) {
                        await this.promisedSetState({groupByColumn: columns.data.items[0]}); 
                    }
                    if (this.state.orderByColumn === null) {
                        await this.promisedSetState({orderByColumn: columns.data.items[0]}); 
                    }
                } else {
                    await this.promisedSetState({columns: []});
                }
            })
            .catch(async() => {
                await this.promisedSetState({columns: []});
            });
    };

    onChangeQuery = async(e) => {
        await this.promisedSetState({
            query: e.target.value,
        });

        this.props.onChangeQuery(e.target.value);
    }

    onChangeQueryFromEditor = async() => {
        let query = this.editorToSQL();
        
        await this.promisedSetState({query});

        this.props.onChangeQuery(query);
    }

    onChangeDatabase = async(value) => {
        await this.promisedSetState({
            database: value,
            tables: null,
        });

        await this.onChangeQueryFromEditor();    

        this.handleGetTables(this.state.integrationUuid, value)
    }

    onChangeTable = async(value) => {
        await this.promisedSetState({
            table: value,
        });

        await this.onChangeQueryFromEditor(); 

        this.handleDescribeTable(
            this.state.integrationUuid,
            this.state.database,
            value
        )
    }

    onChangeFilterColumn = async(value) => {
        await this.promisedSetState({
            filterColumn: value,
        });
    }

    onChangeGroupByColumn = async(value) => {
        await this.promisedSetState({
            groupByColumn: value,
        });
    }

    onChangeOrderByColumn = async(value) => {
        await this.promisedSetState({
            orderByColumn: value,
        });
    }

    onChangeFilterValue = async(e) => {
        await this.promisedSetState({
            filterValue: e.target.value,
        });
    }

    onChangeFilterOperation = async(key) => {
        await this.promisedSetState({
            filterOperation: equationToHuman[key],
            filterOperationSql: equationToSQL[key],
            filterOperationKey: key,
        });
    }

    onChangeOrderDirection = async(key) => {
        await this.promisedSetState({
            orderDirection: orderDirectionToHuman[key],
            orderDirectionSql: orderDirectionToSQL[key],
            orderDirectionKey: key,
        });
    }

    isSQLEditorCompatible = (sql) => {
        const exclusion = new RegExp(/((JOIN|HAVING|COUNT\(\*\)| as|UNION|SHOW)( |\n))/i);
      
        if (sql && exclusion.test(sql)) {
            return false;
        }

        const inclusion = new RegExp(/(SELECT( |\n)+\*( |\n)+FROM( |\n))/i);
      
        if (sql && !inclusion.test(sql)) {
            return false;
        }

        return true;
    }

    render() {
      const { query, showSQLForm, isFilterFormOpen, isGroupByFormOpen, isOrderByFormOpen } = this.state;

      return (
        <>
            <Nav variant="tabs" className="mb-3">
                <Nav.Item>
                    <Nav.Link
                        active={showSQLForm !== true}
                        className="editorTab"
                        onClick={(e) => this.toggleSQLForm(false)}
                    >
                        Constructor
                    </Nav.Link>
                </Nav.Item>
                <Nav.Item>
                    <Nav.Link
                        active={showSQLForm === true}
                        className="editorTab"
                        onClick={(e) => this.toggleSQLForm(true)}
                    >
                        SQL
                    </Nav.Link>
                </Nav.Item>
            </Nav>

            { showSQLForm ?
                <>
                    <CodeEditor
                        className="form-control form-control-lg codeEditor"
                        type="text"
                        language="sql"
                        minHeight={230}
                        autoComplete="SQLQuery"
                        name="SQLQuery"
                        value={query}
                        onChange={this.onChangeQuery}
                        validations={[required]}
                        rehypePlugins={[
                            [rehypePrism, { ignoreMissing: true, showLineNumbers: true }],
                        ]}
                        style={{
                            fontSize: 14,
                            fontFamily: 'Source Code Pro, monospace',
                        }}
                    />
                    <div className="inputTip">
                        Here you can use environment variables: &#123;&#123; ENV_variable_name &#125;&#125;.<br/>
                        And input parameters by their name. For example: &#123;&#123; organization_uuid &#125;&#125;.
                    </div>
                </>
            :
                <div>
                    {
                        !this.isSQLEditorCompatible(query)
                        ?
                            "Your SQL query is complicated and not compatible with visual editor. If you want to use editor, please remove SQL Query."
                        :
                            <div>
                            {
                                this.state.databases !== null 
                                ?
                                    this.state.databases.length > 0
                                    ?
                                        <div className="pb-3">
                                        <span className="inputLabel">Database</span><br/>
                                        <Dropdown size="lg" className="mb-4">
                                            <Dropdown.Toggle 
                                                variant="light" 
                                                id="dropdown-basic"
                                                className="mt-2"
                                            >
                                                <StorageRounded className="sidebarIcon editorIcon pt-1"/>
                                                {this.state.database !== null 
                                                    ? this.state.database
                                                    : this.state.databases !== null && this.state.databases.length > 0 
                                                        ? this.state.databases[0]
                                                        : "Cannot read databases from this integration"
                                                }
                                            </Dropdown.Toggle>

                                            <Dropdown.Menu>
                                                {this.state.databases.map((value, index) => (
                                                    <Dropdown.Item
                                                        value={value}
                                                        key={"db" + index}
                                                        active={value === this.state.database}
                                                        onClick={(e) => this.onChangeDatabase(value)}
                                                    >
                                                        <StorageRounded className="sidebarIcon editorIcon pt-1"/>
                                                        {value}
                                                    </Dropdown.Item>
                                                ))}
                                            </Dropdown.Menu>
                                        </Dropdown>
                                        </div>
                                    : <div className="pt-2 pb-2">No databases were found in this integration</div>
                                : <div className="text-center"><Spinner animation="grow" className="spinner-primary"/></div>
                            }

                        {
                            this.state.tables !== null 
                            ?
                                this.state.tables.length > 0
                                ?
                                    <div className="pb-3">
                                    <span className="inputLabel">Table</span><br/>
                                    <Dropdown size="lg" className="mb-4">
                                        <Dropdown.Toggle 
                                            variant="light" 
                                            id="dropdown-basic"
                                            className="mt-2"
                                        >
                                            <TableChartRounded className="sidebarIcon editorIcon pt-1"/>
                                            {this.state.table !== null 
                                                ? this.state.table
                                                : this.state.tables !== null && this.state.tables.length > 0 
                                                    ? this.state.tables[0]
                                                    : "Cannot read tables from this database"
                                            }
                                        </Dropdown.Toggle>

                                        <Dropdown.Menu>
                                            {this.state.tables.map((value, index) => (
                                                <Dropdown.Item
                                                    value={value}
                                                    key={"tbl" + index}
                                                    active={value === this.state.tables}
                                                    onClick={(e) => this.onChangeTable(value)}
                                                >
                                                    <TableChartRounded className="sidebarIcon editorIcon pt-1"/>
                                                    {value}
                                                </Dropdown.Item>
                                            ))}
                                        </Dropdown.Menu>
                                    </Dropdown>
                                    </div>  
                                : <div className="pt-2 pb-2">No tables were found in this database</div>
                            : this.state.databases !== null  && this.state.databases.length > 0  && <div className="text-center"><Spinner animation="grow" className="spinner-primary"/></div>
                        }

                        {
                            this.state.tables !== null &&
                            <div>
                                <div className="pb-3">
                                    <span className="inputLabel">Filters</span><br/>
                                    {this.state.humanFilters.map((value, index) => (
                                        <Dropdown 
                                            size="lg" 
                                            className="mb-4"
                                            key={"fltr" + index}
                                        >
                                            <Dropdown.Toggle 
                                                variant="light" 
                                                id="dropdown-basic"
                                                className="queryFilter nonEditableFilter mx-2 mt-2"
                                            >
                                                {value}
                                                <div className="filterIconContainer">
                                                    <DeleteOutlineRounded onClick={() => this.removeFilter(index)} className="pb-1 filterIcon"/>
                                                </div>
                                            </Dropdown.Toggle>
                                        </Dropdown>
                                    ))}

                                    {isFilterFormOpen === true ?
                                        <div className="filterForm">
                                            <ClearRounded onClick={() => this.toggleFilterForm()} className="closeFilterFormIcon"/>
                                            { this.state.columns !== null && this.state.columns.length > 0 &&
                                            <Dropdown size="lg" className="mb-4">
                                                <Dropdown.Toggle 
                                                    variant="light" 
                                                    id="dropdown-basic"
                                                    className="mt-2"
                                                >
                                                    {
                                                        this.state.filterColumn !== null
                                                        ? this.state.filterColumn
                                                        : this.state.columns[0]
                                                    }
                                                </Dropdown.Toggle>

                                                <Dropdown.Menu>
                                                    {this.state.columns.map((value, index) => (
                                                        <Dropdown.Item
                                                            value={value}
                                                            key={"col" + index}
                                                            active={value === this.state.filterColumn}
                                                            onClick={(e) => this.onChangeFilterColumn(value)}
                                                        >
                                                            {value}
                                                        </Dropdown.Item>
                                                    ))}
                                                </Dropdown.Menu>
                                            </Dropdown>
                                        }
                                            <Dropdown size="lg" className="mb-4 block">
                                                <Dropdown.Toggle 
                                                    variant="light" 
                                                    id="dropdown-basic"
                                                    className="mt-2 block"
                                                >
                                                    {
                                                        this.state.filterOperation !== null
                                                        ? this.state.filterOperation
                                                        : equationToHuman.neq
                                                    }
                                                </Dropdown.Toggle>

                                                <Dropdown.Menu>
                                                    {Object.keys(equationToHuman).map((key, index) => (
                                                        <Dropdown.Item
                                                            value={equationToHuman[key]}
                                                            key={"col" + index}
                                                            active={equationToHuman[key] === this.state.filterOperation}
                                                            onClick={(e) => this.onChangeFilterOperation(key)}
                                                        >
                                                            {equationToHuman[key]}
                                                        </Dropdown.Item>
                                                    ))}
                                                </Dropdown.Menu>
                                            </Dropdown>
                                            {
                                                !(this.state.filterOperationKey in equationNoValue) &&
                                                <FloatingLabel controlId="floatingValue block" label="Value">
                                                    <Input
                                                        className="form-control form-control-lg mt-2"
                                                        type="text"
                                                        id="floatingValue"
                                                        placeholder="Value"
                                                        autoComplete="value"
                                                        name="value"
                                                        value={this.state.filterValue}
                                                        onChange={this.onChangeFilterValue}
                                                        validations={[required]}
                                                    />
                                                </FloatingLabel>
                                            }
                                            <Button
                                                className="px-3 mt-3 block btn btn-secondary"
                                                onClick={() => this.addFilter()}
                                            >
                                                <span>Add</span>
                                            </Button>
                                        </div>
                                        :
                                        <Dropdown size="lg" className="mb-4">
                                            <Dropdown.Toggle 
                                                variant="light" 
                                                id="dropdown-basic"
                                                className="queryFilter mx-2 mt-2 px-3"
                                                onClick={() => this.toggleFilterForm()}
                                            >
                                                <AddRounded className="sidebarIcon pt-1"/>
                                            </Dropdown.Toggle>
                                        </Dropdown>
                                    }
                                </div>
                                <div className="pb-3">
                                    <span className="inputLabel">Group By</span><br/>
                                    {this.state.groupBys.map((value, index) => (
                                        <Dropdown 
                                            size="lg" 
                                            className="mb-4"
                                            key={"groupBy-" + index}
                                        >
                                            <Dropdown.Toggle 
                                                variant="light" 
                                                id="dropdown-basic"
                                                className="queryFilter nonEditableFilter mx-2 mt-2"
                                            >
                                                {value}
                                                <div className="filterIconContainer">
                                                    <DeleteOutlineRounded onClick={() => this.removeGroupBy(index)} className="pb-1 filterIcon"/>
                                                </div>
                                            </Dropdown.Toggle>
                                        </Dropdown>
                                    ))}

                                    {isGroupByFormOpen === true ?
                                        <div className="filterForm">
                                            <ClearRounded onClick={() => this.toggleGroupByForm()} className="closeFilterFormIcon"/>
                                            { this.state.columns !== null && this.state.columns.length > 0 &&
                                                <Dropdown size="lg" className="mb-4">
                                                    <Dropdown.Toggle 
                                                        variant="light" 
                                                        id="dropdown-basic"
                                                        className="mt-2"
                                                    >
                                                        {
                                                            this.state.groupByColumn !== null
                                                            ? this.state.groupByColumn
                                                            : this.state.columns[0]
                                                        }
                                                    </Dropdown.Toggle>

                                                    <Dropdown.Menu>
                                                        {this.state.columns.map((value, index) => (
                                                            <Dropdown.Item
                                                                value={value}
                                                                key={"col" + index}
                                                                active={value === this.state.groupByColumn}
                                                                onClick={(e) => this.onChangeGroupByColumn(value)}
                                                            >
                                                                {value}
                                                            </Dropdown.Item>
                                                        ))}
                                                    </Dropdown.Menu>
                                                </Dropdown> 
                                            }
                                            
                                            <Button
                                                className="px-3 mt-3 block btn btn-secondary"
                                                onClick={() => this.addGroupBy()}
                                            >
                                                <span>Add</span>
                                            </Button>
                                        </div>
                                        :
                                        <Dropdown size="lg" className="mb-4">
                                            <Dropdown.Toggle 
                                                variant="light" 
                                                id="dropdown-basic"
                                                className="queryFilter mx-2 mt-2 px-3"
                                                onClick={() => this.toggleGroupByForm()}
                                            >
                                                <AddRounded className="sidebarIcon pt-1"/>
                                            </Dropdown.Toggle>
                                        </Dropdown>
                                    }
                                </div>

                                <div className="pb-3">
                                    <span className="inputLabel">Sort By</span><br/>
                                    {this.state.humanOrderBys.map((value, index) => (
                                        <Dropdown 
                                            size="lg" 
                                            className="mb-4"
                                            key={"orderBy-" + index}
                                        >
                                            <Dropdown.Toggle 
                                                variant="light" 
                                                id="dropdown-basic"
                                                className="queryFilter nonEditableFilter mx-2 mt-2"
                                            >
                                                {value}
                                                <div className="filterIconContainer">
                                                    <DeleteOutlineRounded onClick={() => this.removeOrderBy(index)} className="pb-1 filterIcon"/>
                                                </div>
                                            </Dropdown.Toggle>
                                        </Dropdown>
                                    ))}

                                    {isOrderByFormOpen === true ?
                                        <div className="filterForm">
                                            <ClearRounded onClick={() => this.toggleOrderByForm()} className="closeFilterFormIcon"/>
                                            { this.state.columns !== null && this.state.columns.length > 0 &&
                                                <Dropdown size="lg" className="mb-4">
                                                    <Dropdown.Toggle 
                                                        variant="light" 
                                                        id="dropdown-basic"
                                                        className="mt-2"
                                                    >
                                                        {
                                                            this.state.orderByColumn !== null
                                                            ? this.state.orderByColumn
                                                            : this.state.columns[0]
                                                        }
                                                    </Dropdown.Toggle>

                                                    <Dropdown.Menu>
                                                        {this.state.columns.map((value, index) => (
                                                            <Dropdown.Item
                                                                value={value}
                                                                key={"col" + index}
                                                                active={value === this.state.orderByColumn}
                                                                onClick={(e) => this.onChangeOrderByColumn(value)}
                                                            >
                                                                {value}
                                                            </Dropdown.Item>
                                                        ))}
                                                    </Dropdown.Menu>
                                                </Dropdown> 
                                            }

                                            <Dropdown size="lg" className="mb-4 block">
                                                <Dropdown.Toggle 
                                                    variant="light" 
                                                    id="dropdown-basic"
                                                    className="mt-2 block"
                                                >
                                                    {
                                                        this.state.orderDirection !== null
                                                        ? this.state.orderDirection
                                                        : orderDirectionToHuman.nat
                                                    }
                                                </Dropdown.Toggle>

                                                <Dropdown.Menu>
                                                    {Object.keys(orderDirectionToHuman).map((key, index) => (
                                                        <Dropdown.Item
                                                            value={orderDirectionToHuman[key]}
                                                            key={"col-direction" + index}
                                                            active={orderDirectionToHuman[key] === this.state.orderDirection}
                                                            onClick={(e) => this.onChangeOrderDirection(key)}
                                                        >
                                                            {orderDirectionToHuman[key]}
                                                        </Dropdown.Item>
                                                    ))}
                                                </Dropdown.Menu>
                                            </Dropdown>
                                            
                                            <Button
                                                className="px-3 mt-3 block btn btn-secondary"
                                                onClick={() => this.addOrderBy()}
                                            >
                                                <span>Add</span>
                                            </Button>
                                        </div>
                                        :
                                        <Dropdown size="lg" className="mb-4">
                                            <Dropdown.Toggle 
                                                variant="light" 
                                                id="dropdown-basic"
                                                className="queryFilter mx-2 mt-2 px-3"
                                                onClick={() => this.toggleOrderByForm()}
                                            >
                                                <AddRounded className="sidebarIcon pt-1"/>
                                            </Dropdown.Toggle>
                                        </Dropdown>
                                    }
                                </div>
                            </div>
                        }
                    </div>
                }
                </div>
            }
        </>
      );
  }
};

export default control(QueryUI);
