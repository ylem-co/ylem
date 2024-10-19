import React from 'react';
import { Navigate, useParams, useNavigate } from 'react-router-dom';
import {connect} from "react-redux";
import { DateTime } from "luxon";

import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Card from 'react-bootstrap/Card';
import Spinner from "react-bootstrap/Spinner";
import Button from "react-bootstrap/Button";
import ButtonGroup from "react-bootstrap/ButtonGroup";
import Dropdown from "react-bootstrap/Dropdown";

import Circle from '@mui/icons-material/Circle';
import Edit from '@mui/icons-material/Edit';
import ContentCopyRounded from '@mui/icons-material/ContentCopyRounded';
import BarChartRounded from '@mui/icons-material/BarChartRounded';
import HistoryToggleOffRounded from '@mui/icons-material/HistoryToggleOffRounded';
import DeleteOutlineRounded from '@mui/icons-material/DeleteOutlineRounded';
import ContentPasteOutlined from '@mui/icons-material/ContentPasteOutlined';
import GridView from '@mui/icons-material/GridView';
import FormatListBulletedRounded from '@mui/icons-material/FormatListBulletedRounded';
import ContentCopy from '@mui/icons-material/ContentCopy';

import Tooltip from '@mui/material/Tooltip';

import {PipelinesInfo, MetricsInfo} from "../../../actions/infoTexts";
import InfoModal from "../../../components/modals/infoModal.component";
import FullScreenWithMenusModal from "../../../components/modals/fullScreenWithMenusModal.component";
import ConfirmationModal from "../../../components/modals/confirmationModal.component";
import RightModal from "../../../components/modals/rightModal.component";

import LineChart, { lineValueDataSample, lineChartValueOptions } from "../../../components/charts/lineChart.component";

import { PipelineBreadcrumbs } from "../../../components/pipelines/pipelineBreadcrumbs.component";
import { PipelineTriggers } from "../../../components/pipelines/pipelineTriggers.component";
import PipelineLogs from "../../../components/pipelines/pipelineLogs.component";
import PipelineTabs from "../../../components/pipelines/pipelineTabs.component";
import Pipeline from "../../../components/pipelines/pipeline.component";
import PipelineTemplates from "../../../components/pipelines/pipelineTemplates.component";
import FolderForm from "../../../components/forms/folderForm.component";
import ScheduleForm from "../../../components/forms/scheduleForm.component";
import PipelineStatistic from "../../../components/pipelines/statistic/pipelineStatistic.component";
import {TimeAgo} from "../../../components/timeAgo.component";
import MetricForm from "../../../components/forms/metricForm.component";

import PipelineService, 
    {
        PIPELINE_PAGES,
        PIPELINE_PAGE_STATS,
        PIPELINE_PAGE_DETAILS,
        PIPELINE_PAGE_TRIGGERS,
        PIPELINE_PAGE_PREVIEW, 
        PIPELINE_PAGE_LOGS,
        PIPELINE_ROOT_FOLDER, 
        PIPELINE_TYPE_GENERIC, 
        PIPELINE_TYPE_METRIC
    } from "../../../services/pipeline.service";
import FolderService from "../../../services/folder.service";
import TemplateService from "../../../services/template.service";
import StatService from "../../../services/stat.service";

import BootstrapTable from 'react-bootstrap-table-next';
import ToolkitProvider, { Search } from 'react-bootstrap-table2-toolkit/dist/react-bootstrap-table2-toolkit.min';

import {PERMISSION_LOGGED_IN, validatePermissions} from "../../../actions/pipeline";

const { SearchBar } = Search;

const VIEW_GRID = "grid";
const VIEW_LIST = "list";

const borderColors = [
    'rgb(119, 209, 247)',
    'rgb(95, 232, 191)',
    'rgb(243, 114, 44)',
    'rgb(249, 199, 79)',
    'rgb(160, 227, 109)',
    'rgb(170, 122, 235)',
];

function copyLink(copiedLink){
  navigator.clipboard.writeText(copiedLink);
}

function withParams(Component) {
  return props => <Component {...props} params={useParams()} history={useNavigate()} />;
}

class Pipelines extends React.Component {
    constructor(props) {
        super(props);
        this.onDragStart = this.onDragStart.bind(this);
        this.onDrop = this.onDrop.bind(this);
        this.switchView = this.switchView.bind(this);

        this.state = {
            organization: localStorage.getItem('organization') ? JSON.parse(localStorage.getItem('organization')) : [],
            view: localStorage.getItem(this.props.type + '_view') || VIEW_GRID,
            isInfoOpen: false,
            isFormOpen: false,
            isScheduleFormOpen: false,
            isStatisticOpen: false,
            isFolderFormOpen: false,
            isTemplatesListOpen: false,
            isCopyModalOpen: false,
            pipelines: null,
            folders: null,
            activeItem: null,
            activeFolderItem: null,
            statItem: null,
            folderToRemove: null,
            pipelineToRemove: null,
            pipelineToToggle: null,
            pipelineToSaveAsTemplate: null,
            pipelineToCopy: null,
            page: this.props.params.page || null,
            itemUuid: this.props.params.itemUuid || null,
            folderUuid: this.props.params.folderUuid || PIPELINE_ROOT_FOLDER,
            folder: null,
            parentFolder: null,
            isDraggedNow: null,
            isDraggedType: null,
            isDraggedTo: null,
            errorMessage: null,
            previousPage: null,
            isMoveModalOpen: false,
            isSaveAsTemplateModalOpen: false,
            dateFrom: this.props.params.dateFrom || DateTime.now().toSQLDate() + " 00:00:00",
            dateTo: this.props.params.dateTo || DateTime.now().toSQLDate() + " 23:59:59",
            type: this.props.type,
            typeCapital: this.props.type === PIPELINE_TYPE_GENERIC ? "Pipeline" : "Metric",
            typeCapitalPlural: this.props.type === PIPELINE_TYPE_GENERIC ? "Pipelines" : "Metrics",
            typeText: this.props.type === PIPELINE_TYPE_GENERIC ? "pipeline" : "metric",
            typeTextPlural: this.props.type === PIPELINE_TYPE_GENERIC ? "pipelines" : "metrics",
        };
    }

    componentDidMount = async() => {
        document.title = this.state.typeCapitalPlural;

        await this.handleGetPipelines(
            this.state.organization.uuid,
            this.state.folderUuid !== PIPELINE_ROOT_FOLDER
                ? this.state.folderUuid
                : null,
            this.state.itemUuid
        );

        if (this.state.page === null) {
            await this.handleGetFolders(
                this.state.organization.uuid,
                this.state.folderUuid !== PIPELINE_ROOT_FOLDER
                    ? this.state.folderUuid
                    : null
            );
        } else {
            await this.promisedSetState({folders: []});
        }

        if (this.state.folderUuid !== PIPELINE_ROOT_FOLDER) {
            let folder = await this.handleGetFolder(this.state.folderUuid);
            await this.promisedSetState({folder});
        }

        if (
            this.state.folder !== null
            && this.state.folder.parent_uuid !== ""
        ) {
            let parentFolder = await this.handleGetFolder(
                this.state.folder.parent_uuid
            );
            await this.promisedSetState({parentFolder});
        }

        if (
            Object.values(PIPELINE_PAGES).includes(this.state.page)
            && this.state.itemUuid !== null
        ) {
            var element = this.state.pipelines.find(o => o.uuid === this.state.itemUuid);

            if (element) {
                if (this.state.page === PIPELINE_PAGE_DETAILS) {
                    this.toogleForm(element)
                } else if (this.state.page === PIPELINE_PAGE_STATS) {
                    this.toogleStatistic(element)
                } else if (
                    this.state.page === PIPELINE_PAGE_PREVIEW
                    || this.state.page === PIPELINE_PAGE_TRIGGERS
                    || this.state.page === PIPELINE_PAGE_LOGS
                ) {
                    await this.promisedSetState({
                        activeItem: element,
                    });
                }
            }
        }
    };

    promisedSetState = (newState) => new Promise(resolve => this.setState(newState, resolve));

    UNSAFE_componentWillReceiveProps = async(props) => {
        if (
            !props.params.page
            && props.params.folderUuid !== PIPELINE_ROOT_FOLDER
        ) {
            window.location.reload();
        }
    }

    switchView = async(view) => {
        await localStorage.setItem(
            this.state.type + '_view',  view
        );

        this.setState({view});
    }

    handleGetLastMetricValues = async(uuid) => {
        let stats = null;
        stats = StatService.getLastPipelineValues(uuid);

        let data = await Promise.resolve(stats)
            .then(async(stats) => {
                if (stats.data) {
                    var values = stats.data;

                    values.sort((a,b) => new Date(a.executed_at) - new Date(b.executed_at));

                    let lineChartValueData = JSON.parse(JSON.stringify(lineValueDataSample));
                    let lChartValueOptions = JSON.parse(JSON.stringify(lineChartValueOptions)); 
                    
                    let borderColor = borderColors[Math.floor(Math.random() * borderColors.length)];;
                    for(var i = 0; i < values.length; i++){
                        let valueLabel = values[i].executed_at;

                        lineChartValueData.datasets[0].data.push(Number(values[i].metric_value));
                        lineChartValueData.datasets[0].borderColor = borderColor;
                        lineChartValueData.labels.push(valueLabel);
                    }

                    lChartValueOptions.plugins["tooltip"]["callbacks"]["label"] = function(tooltipItem, data) {
                        return tooltipItem.label + ": " + tooltipItem.raw;
                    };

                    return {
                        "data": lineChartValueData,
                        "options": lChartValueOptions,
                    };
                } else {
                    return {
                        "data": lineValueDataSample,
                        "options": lineChartValueOptions,
                    };
                }
            })
            .catch(() => {
                return {
                    "data": lineValueDataSample,
                    "options": lineChartValueOptions,
                };
            });

        return data;
    };

    handleGetFolders = async(uuid, folderUuid = null) => {
        let folders = this.state.folders;

        if (
            folders === null
            || folders.length === 0
        ) {
            folders = FolderService.getFoldersByOrganizationAndFolder(uuid, folderUuid);

            await Promise.resolve(folders)
                .then(async(folders) => {
                    if (folders.data && folders.data.items !== null) {
                        var items = folders.data.items;

                        items = items.filter(k => k.type === this.state.type);

                        await this.promisedSetState({folders: items});
                    } else {
                        await this.promisedSetState({folders: []});
                    }
                });
        }
    };

    handleGetPipelines = async(uuid, folderUuid = null, pipelineUuid = null) => {
        let pipelines = this.state.pipelines;

        if (
            pipelines === null
            || pipelines.length === 0
        ) {
            if (pipelineUuid === null) {
                pipelines = PipelineService.getPipelinesByOrganizationAndFolder(uuid, folderUuid);
            } else {
                pipelines = PipelineService.getPipeline(pipelineUuid);
            }

            await Promise.resolve(pipelines)
                .then(async(pipelines) => {
                    if (
                        pipelines.data 
                        && 
                        (
                            ( 
                                pipelineUuid === null
                                && pipelines.data.items 
                                && pipelines.data.items !== null
                            ) ||
                            (
                                pipelineUuid !== null
                                && pipelines.data.uuid
                            )
                        )
                    ) {
                        if (pipelineUuid === null) {
                            var items = pipelines.data.items.filter(k => k.type === this.state.type);
                        } else {
                            items = [pipelines.data];
                        }

                        await this.promisedSetState({ pipelines: items });

                        if (this.state.type === PIPELINE_TYPE_GENERIC) {
                            this.handleGetPipelinesPreviews(items);
                        } else {
                            this.handleGetLastMetricsValues(items);
                        }
                    } else {
                        await this.promisedSetState({pipelines: []});
                    }
                });
        }
    };

    handleGetFolder = async(uuid) => {
        let folder = FolderService.getFolder(uuid);

        return await Promise.resolve(folder)
            .then(async(folder) => {
                if (folder.data) {
                    return folder.data
                } else {
                    return null
                }
            });
    };

    handleGetPipelinesPreviews = async(items) => {
        var image = "";
        var src = "";

        for(var i = 0; i < items.length; i++){
            if (items[i].preview !== "") {
                image = PipelineService.getPipelinePreview(items[i].uuid);
                src = await Promise.resolve(image)
                .then(value => {
                    return value;
                });
                items[i].src = src;
                this.setState({pipelines: items});
            }
        }
    };

    handleGetLastMetricsValues = async(items) => {
        for(var i = 0; i < items.length; i++){
            if (!items[i].values_data) {
                let values_data = await this.handleGetLastMetricValues(items[i].uuid);

                items[i].values_data = values_data;
                this.setState({pipelines: items});
            }
        }
    };

    toogleInfo = async() => {
        await this.promisedSetState({
            isInfoOpen: !this.state.isInfoOpen,
        });
    };

    closeInfo = () => {
        this.setState({isInfoOpen: false});
    };

    toogleTemplatesList = async() => {
        await this.promisedSetState({
            isTemplatesListOpen: !this.state.isTemplatesListOpen,
        });
    };

    closeTemplatesList = () => {
        this.setState({isTemplatesListOpen: false});
    };

    toogleFolderForm = async(item = null) => {
        await this.promisedSetState({
            isFolderFormOpen: !this.state.isFolderFormOpen,
            activeFolderItem: item,
        });
    };

    closeFolderForm = () => {
        this.setState({
            isFolderFormOpen: false,
            activeFolderItem: null,
        });
    };

    toogleForm = async(item = null, previousPage = PIPELINE_PAGE_PREVIEW) => {
        await this.promisedSetState({
            isFormOpen: !this.state.isFormOpen,
            activeItem: item,
            previousPage,
        });
    };

    closeForm = async(itemUuid = null) => {
        let item = this.state.activeItem;

        if (this.state.type === PIPELINE_TYPE_METRIC) {
            if (item === null && itemUuid !== null) {
                await this.promisedSetState({
                    isFormOpen: false,
                    activeItem: null,
                });
                this.props.history(
                    '/' + this.state.typeTextPlural + '/folder/' + this.state.folderUuid
                    + '/' + PIPELINE_PAGE_DETAILS + '/' + itemUuid
                );
            }
        } else {
            await this.promisedSetState({
                isFormOpen: false,
                activeItem: null,
            });
            if (item !== null) {
                this.props.history(
                    '/' + this.state.typeTextPlural + '/folder/' + this.state.folderUuid
                    + '/' + this.state.previousPage + '/' + item.uuid
                );
                window.location.reload();
            } else {
                this.props.history(
                    '/' + this.state.typeTextPlural + '/folder/' + this.state.folderUuid
                );
            }
        }
    };

    openPreview = (item) => {
        this.props.history(
            '/' + this.state.typeTextPlural + '/folder/' + this.state.folderUuid
            + '/' + PIPELINE_PAGE_PREVIEW + '/' + item.uuid
        );
        window.location.reload();
    }

    toogleScheduleForm = async(item = null, previousPage = null) => {
        await this.promisedSetState({
            isScheduleFormOpen: !this.state.isScheduleFormOpen,
            activeItem: item,
            previousPage,
        });
    };

    closeScheduleForm = async() => {
        let item = this.state.activeItem;

        await this.promisedSetState({
            isScheduleFormOpen: false,
            activeItem: null,
        });
        if (item !== null) {
            if (this.state.previousPage === null) {
                this.props.history('/' + this.state.typeTextPlural + '/folder/' + this.state.folderUuid);
            } else {
                this.props.history(
                    '/' + this.state.typeTextPlural + '/folder/' + this.state.folderUuid
                    + '/' + this.state.previousPage + '/' + item.uuid
                );
                window.location.reload();
            }
        }
    };

    changeDatesHandler = async(dateFrom, dateTo) => {
        await this.promisedSetState({
            dateFrom: dateFrom + '%2000:00:00',
            dateTo: dateTo + '%2023:59:59', 
        });
        this.props.history(
            '/' + this.state.typeTextPlural 
            + '/folder/' + this.state.folderUuid 
            + '/' + this.state.page + '/' + this.state.itemUuid 
            + '/' + dateFrom + '%2000:00:00' 
            + '/' + dateTo + '%2023:59:59');
    };

    openStatistic = (item) => {
        this.toogleStatistic(item);
        this.props.history('/' + this.state.typeTextPlural + '/folder/' + this.state.folderUuid + '/stats/' + item.uuid + '/' + this.state.dateFrom + '/' + this.state.dateTo);
    };

    openForm = (item, previousPage = PIPELINE_PAGE_PREVIEW) => {
        this.toogleForm(item, previousPage);
        
        //if (item.type === PIPELINE_TYPE_METRIC) {
            this.props.history('/' + this.state.typeTextPlural + '/folder/' + this.state.folderUuid + '/details/' + item.uuid);
        //}
    };

    toogleStatistic = async(item = null) => {
        await this.promisedSetState({
            isStatisticOpen: !this.state.isStatisticOpen,
            statItem: item,
        });
    };

    closeStatistic = () => {
        this.props.history('/' + this.state.typeTextPlural + '/folder/' + this.state.folderUuid);
        /*this.setState({
            isStatisticOpen: false,
            statItem: null,
        });*/
    };

    handleSetActiveItem = async(item) => {
        await this.promisedSetState({
            activeItem: item,
            //pipelines: null,
        });
        
        //this.handleGetPipelines(this.state.organization.uuid);
    };

    handleCloseTerminationModal = () => {
        this.setState({
            isTerminationModalOpen: false,
            pipelineToRemove: null,
        });
    };

    handleOpenTerminationModal = (pipeline) => {
        this.setState({
            isTerminationModalOpen: true,
            pipelineToRemove: pipeline,
        });
    };

    handleCloseFolderTerminationModal = () => {
        this.setState({
            isFolderTerminationModalOpen: false,
            folderToRemove: null,
        });
    };

    handleOpenFolderTerminationModal = (folder) => {
        this.setState({
            isFolderTerminationModalOpen: true,
            folderToRemove: folder,
        });
    };

    handleOpenCopyModal = (pipeline) => {
        this.setState({
            isCopyModalOpen: true,
            pipelineToCopy: pipeline,
        });
    };

    handleCloseCopyModal = () => {
        this.setState({
            isCopyModalOpen: false,
            pipelineToCopy: null,
        });
    };

    handleCloseSaveAsTemplateModal = () => {
        this.setState({
            isSaveAsTemplateModalOpen: false,
            pipelineToSaveAsTemplate: null,
        });
    };

    handleOpenSaveAsTemplateModal = (pipeline) => {
        this.setState({
            isSaveAsTemplateModalOpen: true,
            pipelineToSaveAsTemplate: pipeline,
        });
    };

    handleCloseToggleModal = () => {
        this.setState({
            isToggleModalOpen: false,
            pipelineToToggle: null,
        });
    };

    handleOpenToggleModal = (pipeline) => {
        this.setState({
            isToggleModalOpen: true,
            pipelineToToggle: pipeline,
        });
    };

    handleCloseMoveModal = () => {
        this.setState({
            isDraggedNow: null,
            isDraggedType: null,
            isDraggedTo: null,
            isMoveModalOpen: false,
        });
    };

    handleConfirmMove = async() => {
        if (this.state.isDraggedType === "pipeline") {
            await PipelineService.updatePipeline(
                this.state.isDraggedNow.uuid, 
                this.state.isDraggedNow.name, 
                this.state.isDraggedTo !== PIPELINE_ROOT_FOLDER
                    ? this.state.isDraggedTo.uuid
                    : null,
                this.state.isDraggedNow.type === PIPELINE_TYPE_GENERIC
                    ? JSON.parse(this.state.isDraggedNow.elements_layout)
                    : "",
                this.state.isDraggedNow.schedule
            );
            await this.promisedSetState({
                pipelines: null,
            });
            await this.handleGetPipelines(
                this.state.organization.uuid,
                this.state.folderUuid !== PIPELINE_ROOT_FOLDER
                    ? this.state.folderUuid
                    : null
            );
        }

        if (this.state.isDraggedType === "folder") {
            await FolderService.updateFolder(
                this.state.isDraggedNow.uuid, 
                this.state.isDraggedNow.name, 
                this.state.isDraggedTo !== PIPELINE_ROOT_FOLDER
                    ? this.state.isDraggedTo.uuid
                    : ""
            );
            await this.promisedSetState({
                folders: null,
            });
            await this.handleGetFolders(
                this.state.organization.uuid,
                this.state.folderUuid !== PIPELINE_ROOT_FOLDER
                    ? this.state.folderUuid
                    : null
            );
        }

        this.handleCloseMoveModal();
    }

    handleConfirmTermination = async() => {
        var uuid = this.state.pipelineToRemove.uuid;

        await PipelineService.deletePipeline(uuid);
        await this.promisedSetState({
            pipelines: null,
        });
        await this.handleGetPipelines(
            this.state.organization.uuid,
            this.state.folderUuid !== PIPELINE_ROOT_FOLDER
                ? this.state.folderUuid
                : null
        );
        await this.handleCloseTerminationModal();
        await this.closeForm();

        //window.location.reload();
    };

    handleConfirmCopy = async() => {
        var uuid = this.state.pipelineToCopy.uuid;

        let pipeline = await PipelineService.copyPipeline(uuid);

        await Promise.resolve(pipeline)
            .then(async(pipeline) => {
                if(
                    pipeline.uuid
                ) {
                    await this.promisedSetState({
                        errorMessage: null,
                    });
                    await this.handleCloseCopyModal();
                    await this.closeForm();

                    window.location.reload();
                } else {
                    await this.promisedSetState({
                        errorMessage: pipeline.response.data.message,
                    });
                    await this.handleCloseCopyModal();
                    await this.closeForm();
                }
            });
    };

    handleConfirmFolderTermination = async() => {
        var uuid = this.state.folderToRemove.uuid;

        await FolderService.deleteFolder(uuid);
        await this.promisedSetState({
            folders: null,
        });
        await this.handleGetFolders(
            this.state.organization.uuid,
            this.state.folderUuid !== PIPELINE_ROOT_FOLDER
                ? this.state.folderUuid
                : null
        );
        await this.handleCloseFolderTerminationModal();
    };

    handleConfirmToggle = async() => {
        var uuid = this.state.pipelineToToggle.uuid;

        await PipelineService.togglePipeline(uuid);
        await this.promisedSetState({
            pipelines: null,
        });
        await this.handleGetPipelines(
            this.state.organization.uuid,
            this.state.folderUuid !== PIPELINE_ROOT_FOLDER
                ? this.state.folderUuid
                : null
        );
        await this.handleCloseToggleModal();
        await this.closeForm();

        //window.location.reload();
    };

    handleConfirmSaveAsTemplate = async() => {
        var uuid = this.state.pipelineToSaveAsTemplate.uuid;

        await TemplateService.saveAsTemplate(uuid);     
        await this.handleCloseSaveAsTemplateModal();

        window.location.reload();
    };

    handleCloseFormAfterSuccess = async() => {
        await this.promisedSetState({
            pipelines: null,
        });
        await this.handleGetPipelines(
            this.state.organization.uuid,
            this.state.folderUuid !== PIPELINE_ROOT_FOLDER
                ? this.state.folderUuid
                : null
        );
        await this.closeForm();
        
        //window.location.reload();
    }; 

    handleFolderFormSuccess = async() => {
        await this.promisedSetState({
            folders: null,
        });
        await this.handleGetFolders(
            this.state.organization.uuid,
            this.state.folderUuid !== PIPELINE_ROOT_FOLDER
                ? this.state.folderUuid
                : null
        );
        await this.closeFolderForm()
    };

    handleScheduleFormSuccess = async(newSchedule) => {
        let activeItem = this.state.activeItem;
        activeItem.schedule = newSchedule;

        await this.promisedSetState({activeItem});

        await this.closeScheduleForm()
    };

    redirectToFolder = (folder) => {
        this.props.history('/' + this.state.typeTextPlural + '/folder/' + folder.uuid)
    }

    onDragStart = (e, value, type) => {
        e.target.classList.add('draggableActive');
        e.dataTransfer.setData("id", value.uuid);
        Array.from(document.getElementsByClassName('droppable')).forEach(element => (
            element.classList.add('droppableActive')
        ));
        this.setState({isDraggedNow: value, isDraggedType: type})
    };

    onDragOver = (e) => {
        e.preventDefault();
    };

    onDragEnd = (e) => {
        e.target.classList.remove('draggableActive');
        Array.from(document.getElementsByClassName('droppable')).forEach(element => (
            element.classList.remove('droppableActive')
        ))
    };

    onDrop = async(e, value) => {
        if (value.uuid === this.state.isDraggedNow.uuid) {
            return;
        }

        Array.from(document.getElementsByClassName('droppable')).forEach(element => (
            element.classList.remove('droppableActive')
        ));

        this.setState({
            isDraggedTo: value,
            isMoveModalOpen: true
        })
    };

    render() {
        const { 
            isInfoOpen, 
            isFormOpen,
            isTemplatesListOpen,
            isScheduleFormOpen,
            isTerminationModalOpen, 
            isFolderTerminationModalOpen,
            isToggleModalOpen,
            isFolderFormOpen,
            isMoveModalOpen,
            isSaveAsTemplateModalOpen,
            isCopyModalOpen,
            activeItem,
            activeFolderItem,
            statItem,
            isStatisticOpen,
            isDraggedNow,
            isDraggedType,
            isDraggedTo,
            dateFrom,
            dateTo,
            pipelineToToggle,
            errorMessage,
            type,
            folder,
            page,
        } = this.state;

        const { isLoggedIn, user } = this.props;

        if (!validatePermissions(isLoggedIn, user, PERMISSION_LOGGED_IN)) {
            return <Navigate to="/login" />;
        }

        const pipelineColumns = [
            {
                dataField: 'schedule',
                text: '',
                sort: false,
                searchable: false,
                headerStyle: { width: '50px' },
                formatter: (cellContent, row) => (
                    <>
                        <div className="pt-2">
                            <Tooltip title={row.schedule !== ""  ? "Metric is scheduled" : "Metric is not scheduled"} placement="right">
                                <Circle 
                                    className={row.schedule !== "" ? "icon_online" : "icon_new"}
                                    alt={row.schedule !== "" ? "Metric is scheduled" : "Metric is not scheduled"}
                                />
                            </Tooltip>
                        </div>
                    </>
                )
            },
            {
                dataField: 'name',
                text: 'Name',
                sort: false,
                searchable: true,
                formatter: (cellContent, row) => (
                    <>
                        <div 
                            onClick={() => this.openPreview(row)} 
                            className="pt-2 draggable"
                            draggable="true"
                            onDragStart={(e)=>this.onDragStart(e, row, "pipeline")}
                            onDragEnd={(e)=>this.onDragEnd(e)}
                        >
                            {row.name}
                            {
                                row.is_template === 1 &&
                                    <Tooltip title="Template" placement="right">
                                        <span className="templateMarkShort mx-2">T</span>
                                    </Tooltip>
                            }
                        </div>
                    </>
                )
            },
            {
                dataField: 'uuid',
                text: 'UUID',
                sort: false,
                searchable: true,
                formatter: (cellContent, row) => (
                    <>
                        <div onClick={() => this.openPreview(row)} className="pt-2">
                            {row.uuid}
                        </div>
                    </>
                )
            },
            {
                dataField: 'updated_at',
                text: 'Last time updated at',
                sort: false,
                sortFunc: (a, b, order, dataField, rowA, rowB) => {
                    let a1 = rowA.updated_at !== "" 
                        ? DateTime.fromSQL(rowA.updated_at, { zone: 'utc'})
                        : DateTime.fromISO(rowA.created_at, { zone: 'utc'});
                    
                    let b1 = rowB.updated_at !== "" 
                        ? DateTime.fromSQL(rowB.updated_at, { zone: 'utc'})
                        : DateTime.fromISO(rowB.created_at, { zone: 'utc'});

                    if (order === 'asc') {
                      return a1.ts - b1.ts;
                    }
                    return b1.ts - a1.ts;
                  },
                headerStyle: { width: '300px' },
                searchable: true,
                formatter: (cellContent, row) => (
                    <>
                        <div onClick={() => this.openPreview(row)} className="pt-2">
                            {row.updated_at !== "" 
                                ? TimeAgo(DateTime.fromSQL(row.updated_at, { zone: 'utc'})) 
                                : row.created_at !== "" 
                                    ? TimeAgo(DateTime.fromISO(row.created_at, { zone: 'utc'}))
                                    : "New metric"
                            }
                        </div>
                    </>
                )
            },
            {
                dataField: '',
                text: 'Actions',
                sort: false,
                headerStyle: { width: '100px' },
                formatter: (cellContent, row) => (
                    <>
                        <Dropdown as={ButtonGroup}>
                            <Dropdown.Toggle split size="sm" variant="transparent" id="dropdown-split-basic" />

                            <Dropdown.Menu>
                                <Dropdown.Item onClick={() => this.openForm(row)}>
                                    <Edit className="dropdownIcon" /> Edit
                                </Dropdown.Item>
                                <Dropdown.Item onClick={() => this.handleOpenCopyModal(row)}>
                                    <ContentCopyRounded className="dropdownIcon" /> Make a copy
                                </Dropdown.Item>
                                <Dropdown.Item onClick={() => this.openStatistic(row)}>
                                    <BarChartRounded className="dropdownIcon" /> Statistic
                                </Dropdown.Item>
                                <Dropdown.Item onClick={() => this.toogleScheduleForm(row)}>
                                    <HistoryToggleOffRounded className="dropdownIcon" /> Edit schedule
                                </Dropdown.Item>
                                { row.is_template === 0 &&
                                    <Dropdown.Item onClick={() => this.handleOpenSaveAsTemplateModal(row)}>
                                        <ContentPasteOutlined className="dropdownIcon" /> Save as template
                                    </Dropdown.Item>
                                }
                                <Dropdown.Item onClick={() => this.handleOpenTerminationModal(row)} className="dangerAction">
                                    <DeleteOutlineRounded className="dropdownIcon" /> Delete
                                </Dropdown.Item>
                            </Dropdown.Menu>
                        </Dropdown>
                    </>
                )
            },
        ];

        const folderColumns = [
            {
                dataField: 'name',
                text: 'Name',
                sort: false,
                searchable: true,
                formatter: (cellContent, row) => (
                    <>
                        <div 
                            className="folderCard folderInTheList pt-3 droppable onHoverCard"
                            draggable="true"
                            droppable="true"
                            onDragStart={(e)=>this.onDragStart(e, row, "folder")}
                            onDragEnd={(e)=>this.onDragEnd(e)}
                            onDragOver={(e)=>this.onDragOver(e)}
                            onDrop={(e)=>this.onDrop(e, row)}
                        >
                            <a href={'/' + this.state.typeTextPlural + '/folder/' + row.uuid}>
                                {row.name}
                            </a>
                        </div>
                    </>
                )
            },
            {
                dataField: 'updated_at',
                text: 'Last time updated at',
                sort: false,
                sortFunc: (a, b, order, dataField, rowA, rowB) => {
                    let a1 = rowA.updated_at !== "" 
                        ? DateTime.fromSQL(rowA.updated_at, { zone: 'utc'})
                        : DateTime.fromISO(rowA.created_at, { zone: 'utc'});
                    
                    let b1 = rowB.updated_at !== "" 
                        ? DateTime.fromSQL(rowB.updated_at, { zone: 'utc'})
                        : DateTime.fromISO(rowB.created_at, { zone: 'utc'});

                    if (order === 'asc') {
                      return a1.ts - b1.ts;
                    }
                    return b1.ts - a1.ts;
                  },
                headerStyle: { width: '300px' },
                searchable: true,
                formatter: (cellContent, row) => (
                    <>
                        <div className="pt-3">
                            {row.updated_at !== "" 
                                ? TimeAgo(DateTime.fromSQL(row.updated_at, { zone: 'utc'})) 
                                : row.created_at !== "" 
                                    ? TimeAgo(DateTime.fromISO(row.created_at, { zone: 'utc'}))
                                    : "New metric"
                            }
                        </div>
                    </>
                )
            },
            {
                dataField: '',
                text: 'Actions',
                sort: false,
                headerStyle: { width: '100px' },
                formatter: (cellContent, row) => (
                    <>
                        <Dropdown as={ButtonGroup}>
                            <Dropdown.Toggle split size="sm" variant="transparent" id="dropdown-split-basic" />

                            <Dropdown.Menu>
                                <Dropdown.Item onClick={() => this.toogleFolderForm(row)}>
                                    <Edit className="dropdownIcon" /> Edit
                                </Dropdown.Item>
                                <Dropdown.Item className="dangerAction" onClick={() => this.handleOpenFolderTerminationModal(row)}>
                                    <DeleteOutlineRounded className="dropdownIcon" /> Delete
                                </Dropdown.Item>
                            </Dropdown.Menu>
                        </Dropdown>
                    </>
                )
            },
        ];

        const defaultSorted = [{
          dataField: 'updated_at',
          order: 'desc'
        }];

        return (
            <>
                {
                    Object.values(PIPELINE_PAGES).includes(page)
                    && 
                    (
                        <div>
                        {
                            <PipelineBreadcrumbs
                                folder={folder}
                                type={this.state.typeTextPlural}
                                item={
                                    activeItem !== null
                                        ? activeItem
                                        : statItem
                                }
                                parentFolder={this.state.parentFolder}
                                onDrop={this.onDrop}
                                onDragOver={this.onDragOver}
                            />
                        }
                        <PipelineTabs
                            page={page}
                            type={this.state.type}
                            statLink={
                                '/' + this.state.typeTextPlural + '/folder/' + this.state.folderUuid + '/' + PIPELINE_PAGE_STATS + '/' 
                                + this.state.itemUuid
                                + '/' + this.state.dateFrom + '/' + this.state.dateTo
                            }
                            pipelineLink={
                                '/' + this.state.typeTextPlural + '/folder/' + this.state.folderUuid
                                + '/' 
                                + (this.state.type === PIPELINE_TYPE_GENERIC ? PIPELINE_PAGE_PREVIEW : PIPELINE_PAGE_DETAILS)
                                + '/' 
                                + this.state.itemUuid
                                + '/' + this.state.dateFrom + '/' + this.state.dateTo
                            }
                            triggerLink={
                                '/' + this.state.typeTextPlural + '/folder/' + this.state.folderUuid + '/' + PIPELINE_PAGE_TRIGGERS + '/'  
                                + this.state.itemUuid 
                                + '/' + this.state.dateFrom + '/' + this.state.dateTo
                            }
                            logsLink={
                                '/' + this.state.typeTextPlural + '/folder/' + this.state.folderUuid + '/' + PIPELINE_PAGE_LOGS + '/'  
                                + this.state.itemUuid
                            }
                        />
                        </div>
                    )
                }

                {
                    page === PIPELINE_PAGE_TRIGGERS
                    && activeItem !== null
                    &&
                        <PipelineTriggers
                            item={activeItem}
                            openPipelineForm={this.openForm}
                            openScheduleForm={this.toogleScheduleForm}
                        />
                }

                {
                    page === PIPELINE_PAGE_LOGS
                    && activeItem !== null
                    &&
                        <PipelineLogs
                            item={activeItem}
                            dateFrom={dateFrom}
                            dateTo={dateTo}
                            changeDatesHandler={this.changeDatesHandler}
                        />
                }

                { 
                    isFormOpen || page === PIPELINE_PAGE_PREVIEW 
                    ?
                    <div>
                        {
                        type === PIPELINE_TYPE_GENERIC ?
                            <div>
                                {
                                    activeItem !== null
                                    &&
                                        'src' in activeItem
                                            ?
                                                <Row>
                                                    <Col xs={6}>
                                                        {
                                                            activeItem.src !== ""
                                                                ? <img onClick={() => this.openForm(activeItem)} src={'data:image/jpeg;base64,' + activeItem.src} className="imgPreview" alt={"Edit:" + activeItem.name} title={"Edit:" + activeItem.name}/>
                                                                : <div className="pipelinePreview emptyPreview" onClick={() => this.openForm(activeItem)}>Click to edit</div>
                                                        }
                                                    </Col>
                                                    <Col xs={6}>
                                                        <h3>{activeItem.name}</h3>
                                                        {
                                                            activeItem.is_template === 1 &&
                                                            <div>
                                                                <div className="templateMark templateMarkPreview">Template</div>
                                                                <br/><br/>
                                                            </div>
                                                        }
                                                        <div className="code withBorder">
                                                            <Row>
                                                              <Col xs={10}>
                                                                <strong>UUID:</strong> {activeItem.uuid}
                                                              </Col>
                                                              <Col xs={2} className="text-right">
                                                                <Tooltip title="Click to copy to clipboard" placement="left">
                                                                  <ContentCopy
                                                                    onClick={() => copyLink(activeItem.uuid)}
                                                                  />
                                                                </Tooltip>
                                                              </Col>
                                                            </Row>
                                                          </div>
                                                        
                                                        Created at: {DateTime.fromISO(activeItem.created_at, { zone: 'utc'}).toLocal().toLocaleString({ weekday: 'short', month: 'short', year: 'numeric', day: '2-digit', hour: '2-digit', minute: '2-digit' })}<br/>
                                                        {activeItem.updated_at !== "" ? "Last modified: " + TimeAgo(DateTime.fromSQL(activeItem.updated_at, { zone: 'utc'})) : "New " + this.state.typeText }<br/>
                                                        <br/><br/>
                                                        <Button
                                                            variant="primary"
                                                            className="mx-0"
                                                            onClick={() => this.openForm(activeItem)}
                                                          >
                                                            Edit pipeline
                                                          </Button>
                                                    </Col>
                                                </Row>
                                            : <div className="text-center"><Spinner animation="grow" className="spinner-primary"/></div>
                                }
                            </div>
                            :
                            <MetricForm 
                                item={activeItem}
                                successHandler={this.closeForm}
                                folderUuid = {
                                    this.state.folderUuid !== PIPELINE_ROOT_FOLDER
                                        ? this.state.folderUuid
                                        : null
                                }
                              />
                        }
                        
                    </div>
                    :
                        isStatisticOpen
                        ?
                            <PipelineStatistic 
                                item={statItem}
                                dateFrom={dateFrom}
                                dateTo={dateTo}
                                type={type}
                                changeDatesHandler={this.changeDatesHandler}
                            />
                        :
                            !Object.values(PIPELINE_PAGES).includes(page)
                            &&
                    <div>
                        <Row className="mb-2">
                            <Col sm="6">
                                <h1>{this.state.typeCapitalPlural}</h1>
                                <Tooltip title="Info" placement="right">
                                    <div className="infoIcon" onClick={() => this.toogleInfo()}></div>
                                </Tooltip>
                            </Col>
                            <Col sm="6" className="text-right">
                                <ButtonGroup className="mr-4">
                                    {
                                        this.state.view === VIEW_GRID
                                        ? 
                                            <Tooltip title="Switch on list view" placement="left">
                                                <FormatListBulletedRounded
                                                    className="mt-2 mx-4"
                                                    onClick={() => this.switchView(VIEW_LIST)}
                                                />
                                            </Tooltip>
                                        :
                                            <Tooltip title="Switch on grid view" placement="left">
                                                <GridView
                                                    className="mt-2 mx-4"
                                                    onClick={() => this.switchView(VIEW_GRID)}
                                                />
                                            </Tooltip>
                                    }
                                    <Button
                                        variant="secondary"
                                        className="mx-0"
                                        onClick={() => this.toogleFolderForm()}
                                    >
                                        Add folder
                                    </Button>
                                    <Dropdown className="buttonDropdown">
                                        <Dropdown.Toggle className="addPipelineButton btn-primary">
                                            {"Add " + this.state.typeText }
                                        </Dropdown.Toggle>

                                        <Dropdown.Menu>
                                            <Dropdown.Item onClick={() => this.toogleForm()}>Create blank</Dropdown.Item>
                                            <Dropdown.Item onClick={() => this.toogleTemplatesList()}>Create from template</Dropdown.Item>
                                        </Dropdown.Menu>
                                    </Dropdown>
                                </ButtonGroup>
                            </Col>
                            {errorMessage !== null
                                && (
                              <div className="form-group mt-4">
                                <div className="alert alert-danger" role="alert">
                                  {errorMessage}
                                </div>
                              </div>
                            )}
                        </Row>

                {   
                    this.state.pipelines !== null
                    && this.state.folders !== null
                    ?
                    <div>
                        <PipelineBreadcrumbs
                            folder={this.state.folder}
                            type={this.state.typeTextPlural}
                            parentFolder={this.state.parentFolder}
                            onDrop={this.onDrop}
                            onDragOver={this.onDragOver}
                        />
                        {this.state.folders.length > 0 &&
                            <div>
                                { this.state.view === VIEW_GRID
                                    ?
                                <Row>
                                    {this.state.folders.map(value => (
                                        <Col className="col-3 mb-4" key={value.uuid}>
                                            <Card 
                                                className="withHeader onHoverCard folderCard droppable"
                                                draggable="true"
                                                droppable="true"
                                                onDragStart={(e)=>this.onDragStart(e, value, "folder")}
                                                onDragEnd={(e)=>this.onDragEnd(e)}
                                                onDragOver={(e)=>this.onDragOver(e)}
                                                onDrop={(e)=>this.onDrop(e, value)}
                                            >
                                                <Card.Header className="noBottom">
                                                    <Row>
                                                        <Col 
                                                            className="col-7 folderLink"
                                                        >
                                                            <a href={'/' + this.state.typeTextPlural + '/folder/' + value.uuid}>
                                                                {value.name}
                                                            </a>
                                                        </Col>
                                                        <Col className="col-5">
                                                            <div className="text-right cardIcons">
                                                                <Dropdown as={ButtonGroup}>
                                                                    <Dropdown.Toggle split size="sm" variant="transparent" id="dropdown-split-basic" />

                                                                    <Dropdown.Menu>
                                                                        <Dropdown.Item onClick={() => this.toogleFolderForm(value)}>
                                                                            <Edit className="dropdownIcon" /> Edit
                                                                        </Dropdown.Item>
                                                                        <Dropdown.Item className="dangerAction" onClick={() => this.handleOpenFolderTerminationModal(value)}>
                                                                            <DeleteOutlineRounded className="dropdownIcon" /> Delete
                                                                        </Dropdown.Item>
                                                                    </Dropdown.Menu>
                                                                </Dropdown>
                                                            </div> 
                                                        </Col>
                                                    </Row>
                                                </Card.Header>
                                            </Card>
                                        </Col>
                                    ))} 
                                </Row>
                                :
                                    <ToolkitProvider
                                        bootstrap4
                                        keyField="name"
                                        data={ this.state.folders }
                                        columns={ folderColumns }
                                        search
                                    >
                                        {
                                        props => (
                                            <Card className="withHeader mb-5">
                                                {this.state.folders !== null &&
                                                    <>
                                                        <Card.Header className="noBottom">
                                                            <SearchBar {...props.searchProps} srText="" />
                                                        </Card.Header>
                                                        <Card.Body>
                                                            <BootstrapTable
                                                                {...props.baseProps}
                                                                keyField="name"
                                                                bordered={false}
                                                                hover
                                                                defaultSorted={defaultSorted}
                                                                rowClasses="detailedTableRow droppable onHoverCard"
                                                            />
                                                        </Card.Body>
                                                    </>
                                                }
                                            </Card>
                                        )
                                        }
                                    </ToolkitProvider>
                            }
                            </div>
                        } 
                         
                        {
                        this.state.pipelines.length > 0 ?
                        <div>
                            { this.state.view === VIEW_GRID
                                ?
                            <Row>
                            {this.state.pipelines.map(value => (
                                    <Col className="col-3 mb-4" key={value.uuid}>
                                        <Card 
                                            className="withHeader onHoverCard draggable"
                                            draggable="true"
                                            onDragStart={(e)=>this.onDragStart(e, value, "pipeline")}
                                            onDragEnd={(e)=>this.onDragEnd(e)}
                                        >
                                            <Card.Header>
                                                <Row>
                                                    <Col className="col-10">
                                                        {value.name}<br/>
                                                        <span className="undernote">UUID: {value.uuid}</span>
                                                    </Col>
                                                    <Col className="col-2 text-right">
                                                        <Tooltip title={value.schedule !== ""  ? this.state.typeCapital + " is scheduled" : this.state.typeCapital + " is not scheduled"} placement="right">
                                                            <Circle 
                                                                className={value.schedule !== "" ? "icon_online" : "icon_new"}
                                                                alt={value.schedule !== "" ? this.state.typeCapital + " is scheduled" : this.state.typeCapital + " is not scheduled"}
                                                            />
                                                        </Tooltip>
                                                    </Col>
                                                </Row>
                                            </Card.Header> 
                                            <Card.Body>
                                                <div className={
                                                    value.type === PIPELINE_TYPE_GENERIC
                                                        ? "pipelineInTheList mb-3"
                                                        : "mb-3"
                                                    }
                                                >
                                                    {
                                                        value.is_template === 1 &&
                                                        <div className="templateMark">Template</div>
                                                    }
                                                    {
                                                        value.type === PIPELINE_TYPE_GENERIC ?
                                                        (
                                                        'src' in value ?
                                                            value.src !== ""
                                                            ? <img src={'data:image/jpeg;base64,' + value.src} onClick={() => this.openPreview(value)} className="pipelinePreview" alt={value.name} title={value.name}/>
                                                            : <div className="pipelinePreview" onClick={() => this.openPreview(value)}></div>
                                                        : <div className="text-center"><Spinner animation="grow" className="spinner-primary"/></div>
                                                        )
                                                        :
                                                        (
                                                        'values_data' in value ?
                                                            <div className="pipelinePreview" onClick={() => this.openForm(value)}>
                                                                <LineChart
                                                                    data={value.values_data.data}
                                                                    options={value.values_data.options}
                                                                />
                                                            </div>
                                                        : <div className="text-center"><Spinner animation="grow" className="spinner-primary"/></div>
                                                        )
                                                    }
                                                </div>
                                                <span className="note">{value.updated_at !== "" ? "Last modified: " + TimeAgo(DateTime.fromSQL(value.updated_at, { zone: 'utc'})) : "New " + this.state.typeText }</span>
                                                <div className="text-right pipelineCardIcons">
                                                    <Dropdown as={ButtonGroup}>
                                                        <Dropdown.Toggle split size="sm" variant="transparent" id="dropdown-split-basic" />

                                                        <Dropdown.Menu>
                                                            <Dropdown.Item onClick={() => this.openForm(value)}>
                                                                <Edit className="dropdownIcon" /> Edit
                                                            </Dropdown.Item>
                                                            <Dropdown.Item onClick={() => this.handleOpenCopyModal(value)}>
                                                                <ContentCopyRounded className="dropdownIcon" /> Make a copy
                                                            </Dropdown.Item>
                                                            <Dropdown.Item onClick={() => this.openStatistic(value)}>
                                                                <BarChartRounded className="dropdownIcon" /> Statistic
                                                            </Dropdown.Item>
                                                            <Dropdown.Item onClick={() => this.toogleScheduleForm(value)}>
                                                                <HistoryToggleOffRounded className="dropdownIcon" /> Edit schedule
                                                            </Dropdown.Item>
                                                            { value.is_template === 0 &&
                                                                <Dropdown.Item onClick={() => this.handleOpenSaveAsTemplateModal(value)}>
                                                                    <ContentPasteOutlined className="dropdownIcon" /> Save as template
                                                                </Dropdown.Item>
                                                            }
                                                            <Dropdown.Item onClick={() => this.handleOpenTerminationModal(value)} className="dangerAction">
                                                                <DeleteOutlineRounded className="dropdownIcon" /> Delete
                                                            </Dropdown.Item>
                                                        </Dropdown.Menu>
                                                    </Dropdown>
                                                </div> 
                                            </Card.Body>
                                        </Card>
                                    </Col>
                        ))} 
                            </Row>
                            :
                                <ToolkitProvider
                                        bootstrap4
                                        keyField="name"
                                        data={ this.state.pipelines }
                                        columns={ pipelineColumns }
                                        search
                                    >
                                        {
                                        props => (
                                            <Card className="withHeader">
                                                {this.state.pipelines !== null &&
                                                    <>
                                                        <Card.Header className="noBottom">
                                                            <SearchBar {...props.searchProps} srText="" />
                                                        </Card.Header>
                                                        <Card.Body>
                                                            <BootstrapTable
                                                                {...props.baseProps}
                                                                keyField="name"
                                                                bordered={false}
                                                                hover
                                                                defaultSorted={defaultSorted}
                                                                rowClasses="detailedTableRow draggable"
                                                            />
                                                        </Card.Body>
                                                    </>
                                                }
                                            </Card>
                                        )
                                        }
                                    </ToolkitProvider>
                        }
                        </div>
                        : <div className="text-center">
                            {type === PIPELINE_TYPE_GENERIC 
                                ? "You didn't create any pipeline in this folder yet" 
                                : "You didn't create any metric in this folder yet"
                            }
                            </div>
                        }
                    </div>
                : <div className="text-center"><Spinner animation="grow" className="spinner-primary"/></div>
                }
                </div>

            }

                <InfoModal
                    show={isInfoOpen}
                    onHide={this.closeInfo}
                    title={type === PIPELINE_TYPE_GENERIC ? PipelinesInfo.title : MetricsInfo.title}
                    content={type === PIPELINE_TYPE_GENERIC ? PipelinesInfo.content : MetricsInfo.content}
                />

                <InfoModal
                    show={isTemplatesListOpen}
                    onHide={this.closeTemplatesList}
                    title={"Create " + this.state.typeText + " from template"}
                    content={<>
                        <PipelineTemplates
                            folderUuid={this.state.folderUuid}
                            type={type}
                        />
                    </>}
                />

                <FullScreenWithMenusModal
                    show={isFormOpen && this.state.type === PIPELINE_TYPE_GENERIC}
                    onHide={this.handleCloseFormAfterSuccess}
                    title=""
                    content={
                        type === PIPELINE_TYPE_GENERIC ?
                            <Pipeline
                                handleSetActiveItem={this.handleSetActiveItem}
                                item={activeItem}
                                elements={activeItem !== null && activeItem.elements_layout !== "" ? activeItem.elements_layout : "[]"}
                                folderUuid = {
                                    this.state.folderUuid !== PIPELINE_ROOT_FOLDER
                                        ? this.state.folderUuid
                                        : null
                                }
                                newTitle={
                                    activeItem !== null && activeItem.name 
                                        ? activeItem.name 
                                        : this.state.pipelines !== null
                                            ? "New " + this.state.typeText + " " + (this.state.pipelines.length + 1)
                                            : "New " + this.state.typeText 
                                    }
                            />
                            :
                            <MetricForm 
                                item={activeItem}
                                successHandler={this.closeForm}
                                folderUuid = {
                                    this.state.folderUuid !== PIPELINE_ROOT_FOLDER
                                        ? this.state.folderUuid
                                        : null
                                }
                              />
                        
                    }
                    altButton={false}
                    altButton2={false}
                    item={activeItem}
                    confirmButton={false}
                />

                {/*
                <FullScreenWithMenusModal
                    show={isStatisticOpen}
                    onHide={this.closeStatistic}
                    title={statItem !== null && statItem.name ? statItem.name + ". Statistic" :  this.state.typeCapital + " statistic"}
                    content={
                        <>
                            <PipelineStatistic 
                                item={statItem}
                                dateFrom={dateFrom}
                                dateTo={dateTo}
                                type={type}
                                changeDatesHandler={this.changeDatesHandler}
                            />
                        </>
                    }
                    altButton={false}
                    item={statItem}
                    confirmButton={false}
                />
                */}

                <RightModal
                    show={isFolderFormOpen}
                    content={
                      <FolderForm 
                        item={activeFolderItem}
                        successHandler={this.handleFolderFormSuccess}
                        type={type}
                        parentFolder = {
                            this.state.folderUuid !== PIPELINE_ROOT_FOLDER
                                ? this.state.folderUuid
                                : null
                        }
                      />
                    }
                    item={activeFolderItem}
                    title={activeFolderItem === null ? "Add folder" : "Edit folder"}
                    onHide={this.closeFolderForm}
                  />

                <RightModal
                    show={isScheduleFormOpen}
                    content={
                      <ScheduleForm 
                        item={activeItem}
                        successHandler={this.handleScheduleFormSuccess}
                        type={this.state.type}
                      />
                    }
                    item={activeItem}
                    title={"Edit " + this.state.typeText + " running schedule"}
                    onHide={this.closeScheduleForm}
                  />

                <ConfirmationModal
                    show={isFolderTerminationModalOpen}
                    title="Delete folder"
                    body={"Are you sure you want to delete this folder with its content?"}
                    confirmText="Delete"
                    onCancel={this.handleCloseFolderTerminationModal}
                    onHide={this.handleCloseFolderTerminationModal}
                    onConfirm={this.handleConfirmFolderTermination}
                />

                <ConfirmationModal
                    show={isSaveAsTemplateModalOpen}
                    title={"Save " + this.state.typeText + " as template"}
                    body={"Are you sure you want to save this " + this.state.typeText + " as a template?"}
                    confirmText="Save"
                    onCancel={this.handleCloseSaveAsTemplateModal}
                    onHide={this.handleCloseSaveAsTemplateModal}
                    onConfirm={this.handleConfirmSaveAsTemplate}
                />

                <ConfirmationModal
                    show={isTerminationModalOpen}
                    title={"Delete " + this.state.typeText}
                    body={"Are you sure you want to delete this " + this.state.typeText + "?"}
                    confirmText="Delete"
                    onCancel={this.handleCloseTerminationModal}
                    onHide={this.handleCloseTerminationModal}
                    onConfirm={this.handleConfirmTermination}
                />

                <ConfirmationModal
                    show={isCopyModalOpen}
                    title={"Copy " + this.state.typeText}
                    body={"Would you like to make a copy of this " + this.state.typeText + "?"}
                    confirmText="Copy"
                    onCancel={this.handleCloseCopyModal}
                    onHide={this.handleCloseCopyModal}
                    onConfirm={this.handleConfirmCopy}
                />

                <ConfirmationModal
                    show={isToggleModalOpen}
                    title={pipelineToToggle !== null && pipelineToToggle.is_paused === 0 ? "Pause " + this.state.typeText : "Activate " + this.state.typeText}
                    body={
                        pipelineToToggle !== null && pipelineToToggle.is_paused === 0 
                        ? "Are you sure you want to pause this " + this.state.typeText + "?" 
                        : "Are you sure you want to activate this " + this.state.typeText + "?"
                    }
                    confirmText={pipelineToToggle !== null && pipelineToToggle.is_paused === 0 ? "Pause" : "Activate"}
                    onCancel={this.handleCloseToggleModal}
                    onHide={this.handleCloseToggleModal}
                    onConfirm={this.handleConfirmToggle}
                />

                <ConfirmationModal
                    show={isMoveModalOpen}
                    title={"Move " + isDraggedType}
                    body={
                        isDraggedNow !== null && isDraggedTo !== null
                        &&
                        (isDraggedTo !== PIPELINE_ROOT_FOLDER
                            ? 'Are you sure you want to move "' + isDraggedNow.name + '" to the folder "' + isDraggedTo.name + '"?'
                            : 'Are you sure you want to move "' + isDraggedNow.name + '" to the root folder?'
                        )}
                    confirmText={"Move"}
                    onCancel={this.handleCloseMoveModal}
                    onHide={this.handleCloseMoveModal}
                    onConfirm={this.handleConfirmMove}
                />
            </>
        )
    }
}

function mapStateToProps(state) {
    const { isLoggedIn } = state.auth;
    const { user } = state.auth;
    return {
        isLoggedIn,
        user,
    };
}

export default connect(mapStateToProps)(withParams(Pipelines));
