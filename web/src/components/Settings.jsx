import React, { useState } from 'react';
import { Layout, Menu, Button, Modal, theme, Breadcrumb, Space, Form, Input, Checkbox } from 'antd';
import { LaptopOutlined, NotificationOutlined, UserOutlined } from '@ant-design/icons';
import * as reddit from '../lib/reddit';

const { Content, Sider } = Layout;



const onFinish = values => {
    console.log('Success:', values);
};
const onFinishFailed = errorInfo => {
    console.log('Failed:', errorInfo);
};

const redditFlow = async () => {
    const params = new URLSearchParams(window.location.search);

    const code = params.get('code');
    const state = params.get('state');

    if (state === reddit.state && code) {
        const [err, data] = await reddit.requestToken(code);

        console.log(err, data);
        // TODO verify code with request?
        // if passes save to localStorage?

        // TODO get access token through server
        // request with code ...
    }
};

redditFlow();

function Settings() {

    const {
        token: { colorBgContainer, borderRadiusLG },
    } = theme.useToken();

    // TODO if already setup change button to outlined? danger?
    // TODO modal should hook up to backend

    return (
        <>
            <Breadcrumb
                items={[{ title: 'Settings' }]}
            />
            <Layout
                style={{ padding: '24px 0', background: colorBgContainer, borderRadius: borderRadiusLG }}
            >
                <Content style={{ padding: '0 24px', minHeight: 280 }}>
                    <Space>
                        <label>
                            <i>
                            </i>
                            Reddit Connection
                        </label>
                        <Button type="primary" href="/api/v1/reddit/connect" >Setup</Button>
                    </ Space>
                </Content>
            </Layout>
        </>
    );
}


export default Settings;
