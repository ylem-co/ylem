import React, { useState } from 'react';

import Form from 'react-bootstrap/Form';
import Nav from 'react-bootstrap/Nav';
import Navbar from 'react-bootstrap/Navbar';

import { NavLink } from "react-router-dom";

import SettingsOutlined from '@mui/icons-material/SettingsOutlined';
import GroupsOutlined from '@mui/icons-material/GroupsOutlined';
import BlurOnOutlined from '@mui/icons-material/BlurOnOutlined';
import ExploreOutlined from '@mui/icons-material/ExploreOutlined';
import DeviceHubOutlined from '@mui/icons-material/DeviceHubOutlined';
import IntegrationInstructionsOutlined from '@mui/icons-material/IntegrationInstructionsOutlined';
import Fingerprint from '@mui/icons-material/Fingerprint';
import DonutLargeOutlined from '@mui/icons-material/DonutLargeOutlined';
import ExplicitOutlined from '@mui/icons-material/ExplicitOutlined';
import HandymanOutlined from '@mui/icons-material/HandymanOutlined';

import Tooltip from '@mui/material/Tooltip';

import {ROLE_ORGANIZATION_ADMIN} from "../actions/roles";

const Sidebar = () => {
  const user = JSON.parse(localStorage.getItem('user'));
  const sidebarExpanded = localStorage.getItem('sidebarExpanded');
  const darkTheme = localStorage.getItem('darkTheme') !== "false";
  const activeSidebar = localStorage.getItem('activeSidebar');

  const [isExpanded, setIsExpanded] = useState(sidebarExpanded);
  const [isDarkThemeEnabled, setIsDarkThemeEnabled] = useState(darkTheme);
  const [activeExpandedSidebar, setActiveSidebar] = useState(activeSidebar);

  const expand = (activeSidebar) => {
    setIsExpanded(true);
    setActiveSidebar(activeSidebar);
    localStorage.setItem('sidebarExpanded', true);
    localStorage.setItem('activeSidebar', activeSidebar);
  };

  const collapse = () => {
    setIsExpanded(false);
    setActiveSidebar(false);
    localStorage.removeItem('sidebarExpanded');
    localStorage.removeItem('activeSidebar');
  };

  const toggleDarkTheme = () => {
    if (isDarkThemeEnabled) {
      document.body.classList.remove("darkTheme");
      document.documentElement.setAttribute('data-color-mode', 'light');
      setIsDarkThemeEnabled(false);
      localStorage.setItem('darkTheme', false);
      window.location.reload();
      return;
    }

    document.body.classList.add("darkTheme");
    document.documentElement.setAttribute('data-color-mode', 'dark');
    setIsDarkThemeEnabled(true);
    localStorage.setItem('darkTheme', true);
    window.location.reload();
  };

  return (
    <Navbar className={
      isExpanded
       ? "h-100 align-items-start sidebar"
       : "h-100 align-items-start sidebar collapsed"
    }>
      <Nav className="h-100 flex-column px-2 pt-2 sidebarInner">
        <div className="nav-link">
          <a href="/dashboard">
            {isDarkThemeEnabled === true
              ? <img src="/images/logo-s.png" width="20px" alt="Ylem" title="Ylem"/>
              : <img src="/images/logo-s-dark.png" width="20px" alt="Ylem" title="Ylem"/>
            }
          </a>
          {/* isExpanded ? 
            <Tooltip title="Close sidebar" placement="right">
              <MenuOpenRounded
                className="sidebarIcon collapseSidebarIcon"
                onClick={handleToggler}
              />
            </Tooltip>
            :
            <Tooltip title="Open sidebar" placement="right">
              <MenuRounded
                className="sidebarIcon collapseSidebarIcon"
                onClick={handleToggler}
              />
            </Tooltip>
          */}
        </div>
        <NavLink to="/dashboard" className="nav-link firstSidebarItem tour-step-dashboard">
            <Tooltip title="Dashboard" placement="right">
              <BlurOnOutlined 
                className="sidebarIcon"
                onMouseOver={collapse}
              />
            </Tooltip>
        </NavLink>
        <NavLink to="/pipelines" className="nav-link tour-step-pipelines">
            <Tooltip title="Pipelines" placement="right">
              <DeviceHubOutlined 
                className="sidebarIcon"
                onMouseOver={collapse}
              />
            </Tooltip>
        </NavLink>
        <NavLink to="/integrations" className="nav-link tour-step-integrations">
            <Tooltip title="Integrations" placement="right">
              <ExploreOutlined 
                className="sidebarIcon"
                onMouseOver={collapse}
              />
            </Tooltip>
        </NavLink>
        <NavLink to="/metrics" className="nav-link tour-step-metrics">
            <Tooltip title="Metrics" placement="right">
              <DonutLargeOutlined 
                className="sidebarIcon"
                onMouseOver={collapse}
              />
            </Tooltip>
        </NavLink>
        <NavLink to="/slow-tasks" className="nav-link tour-step-profiling">
            <Tooltip title="Profiling" placement="right">
              <HandymanOutlined 
                className="sidebarIcon"
                onMouseOver={collapse}
              />
            </Tooltip>
        </NavLink>
        <NavLink to="/slack-authorizations"
          className={
            activeExpandedSidebar === "authorizations"
              ? "nav-link active"
              : "nav-link"
          }
        >
            <Tooltip title="Authorizations" placement="right">
              <IntegrationInstructionsOutlined 
                className="sidebarIcon"
                onMouseOver={()=>expand("authorizations")}
              />
            </Tooltip>
        </NavLink>
        <NavLink to="/env-variables" className="nav-link tour-step-env-variables">
            <Tooltip title="Environment variables" placement="right">
              <ExplicitOutlined 
                className="sidebarIcon"
                onMouseOver={collapse}
              />
            </Tooltip>
        </NavLink>
        {
          user !== null 
          && user.roles 
          && user.roles.includes(ROLE_ORGANIZATION_ADMIN) 
          &&
        <NavLink 
          to="/users"
          className="nav-link"
        >
            <Tooltip title="Users" placement="right">
              <GroupsOutlined 
                className="sidebarIcon"
                onMouseOver={collapse}
              />
            </Tooltip>
        </NavLink>
        }
        <NavLink 
          to="/api-clients"
          className="nav-link tour-step-oauth-clients"
        >
            <Tooltip title="API Clients" placement="right">
              <Fingerprint 
                className="sidebarIcon"
                onMouseOver={collapse}
              />
            </Tooltip>
        </NavLink>
        <NavLink to="/settings" className="nav-link">
            <Tooltip title="Settings" placement="right">
              <SettingsOutlined 
                className="sidebarIcon"
                onMouseOver={collapse}
              />
            </Tooltip>
        </NavLink>
        <Tooltip title={isDarkThemeEnabled === true ? "Light theme" : "Dark theme"} placement="right">
          <Form.Check
            type="switch"
            className="themeSwitch"
            onChange={toggleDarkTheme}
            checked={isDarkThemeEnabled === true}
          />
        </Tooltip>
      </Nav>
      { activeExpandedSidebar === "authorizations"
        &&
        <Nav className="h-100 flex-column px-2 submenu submenuAuthorizations">
          <h3 className="submenuHeader">Authorizations</h3>
          <NavLink to="/slack-authorizations" className="nav-link">
            <span className="sidebarText">Slack</span>
          </NavLink>
          <NavLink to="/jira-authorizations" className="nav-link">
            <span className="sidebarText">Jira Cloud</span>
          </NavLink>
          <NavLink to="/hubspot-authorizations" className="nav-link">
            <span className="sidebarText">Hubspot</span>
          </NavLink>
          <NavLink to="/salesforce-authorizations" className="nav-link">
            <span className="sidebarText">Salesforce</span>
          </NavLink>
        </Nav>
      }
      {/* activeExpandedSidebar === "profiling"
        &&
        <Nav className="h-100 flex-column px-2 submenu submenuProfiling">
          <h3 className="submenuHeader">Profiling</h3>
          <NavLink to="/slow-tasks" className="nav-link">
            <span className="sidebarText">Slow tasks</span>
          </NavLink>
        </Nav>
      */}
    </Navbar>
  )
}

export default Sidebar
