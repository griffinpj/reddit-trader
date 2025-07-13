import React, { useState } from 'react';
import { Layout, Menu, Button, Modal, theme, Row, Col, Statistic } from 'antd';
import { LaptopOutlined, NotificationOutlined, UserOutlined } from '@ant-design/icons';
import CountUp from 'react-countup';
const { Content, Sider } = Layout;

const formatter = value => <CountUp end={value} separator="," />;

function Home() {
    const {
        token: { colorBgContainer, borderRadiusLG },
    } = theme.useToken();

    return (
        <Content style={{ padding: '0 24px', minHeight: 280 }}>
            <Row gutter={16}>
                <Col span={12}>
                    <Statistic title="Active Posts" value={10} formatter={formatter} />
                </Col>
                <Col span={12}>
                    <Statistic title="Sold" value={11} precision={2} formatter={formatter} />
                </Col>
            </Row>
            <br />
            <Row gutter={16}>
                <Col span={12}>
                    <Statistic title="Buying" value={3} formatter={formatter} />
                </Col>
                <Col span={12}>
                    <Statistic title="Selling" value={7} precision={2} formatter={formatter} />
                </Col>
            </Row>
        </Content>
    );
}


export default Home;
