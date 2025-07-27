import React, { useState } from 'react';
import { Routes, Route, useLocation, useNavigate } from 'react-router';
import { Layout, Breadcrumb, Menu, theme, Button } from 'antd';
import {
    MenuFoldOutlined,
    MenuUnfoldOutlined,
    UploadOutlined,
    UserOutlined,
    HomeOutlined,
    VideoCameraOutlined,
    SettingOutlined,
    UnorderedListOutlined,
    LogoutOutlined
} from '@ant-design/icons';

import * as auth from '../lib/auth';
import Posts from './Posts';
import Home from './Home';
import Settings from './Settings';

const { Header, Content, Footer, Sider } = Layout;

function App() {
    const [collapsed, setCollapsed] = useState(true);
    const {
        token: { colorBgContainer, borderRadiusLG },
    } = theme.useToken();

    const location = useLocation();
    const navigate = useNavigate();


    return (
        <Layout style={{ minHeight: '100vh' }}>
            <Sider trigger={null} collapsible collapsed={collapsed}>
                <div className="demo-logo-vertical" />
                <Menu
                    theme="dark"
                    mode="inline"
                    defaultSelectedKeys={[location.pathname]}
                    items={[
                        {
                            key: '/',
                            icon: <HomeOutlined />,
                            label: 'Home',
                            onClick: () => {
                                navigate('/');
                            }
                        },
                        {
                            key: '/posts',
                            icon: <UnorderedListOutlined/>,
                            label: 'Posts',
                            onClick: () => {
                                navigate('/posts');
                            }
                        },
                        {
                            key: '/settings',
                            icon: <SettingOutlined />,
                            label: 'Settings',
                            onClick: () => {
                                navigate('/settings');
                            }
                        },
                        {
                            key: '/logout',
                            icon: <LogoutOutlined />,
                            label: 'Logout',
                            onClick: async () => {
                                let err = await auth.logout();
                                if (!err) {
                                    navigate('/login');
                                }
                            }
                        },
                    ]}
                />
            </Sider>
            <Layout>
                <Header style={{ padding: 0, background: colorBgContainer }}>
                    <Button
                        type="text"
                        icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
                        onClick={() => setCollapsed(!collapsed)}
                        style={{
                            fontSize: '16px',
                            width: 64,
                            height: 64,
                        }}
                    />
                </Header>
                <Content style={{ padding: '48px 48px 0px 48px' }}>
                    <div
                        style={{
                            background: colorBgContainer,
                            minHeight: 280,
                            padding: 24,
                            borderRadius: borderRadiusLG,
                        }}
                    >
                        <Routes>
                            <Route path="/" element={<Home />} />
                            <Route path="/posts" element={<Posts />} />
                            <Route path="/settings" element={<Settings />} />
                        </Routes>
                    </div>
                </Content>
                <Footer style={{ textAlign: 'center' }}>
                    Ant Design ©{new Date().getFullYear()} Created by Griffin
                </Footer>
            </Layout>
        </Layout>
    );
}

export default App;
