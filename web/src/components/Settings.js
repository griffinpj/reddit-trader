import React, { useState } from 'react';
import { Layout, Menu, Button, Modal, theme, Breadcrumb, Space } from 'antd';
import { LaptopOutlined, NotificationOutlined, UserOutlined } from '@ant-design/icons';
const { Content, Sider } = Layout;

function Settings() {
    const [isModalOpen, setIsModalOpen] = useState(false);
    const showModal = () => {
        setIsModalOpen(true);
    };
    const handleOk = () => {
        setIsModalOpen(false);
    };
    const handleCancel = () => {
        setIsModalOpen(false);
    };

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
                        Reddit Connection
                    </label>
                    <Button type="primary" onClick={showModal}>Setup</Button>
                    </ Space>
                    <Modal
                        title="Connection"
                        closable={{ 'aria-label': 'Custom Close Button' }}
                        open={isModalOpen}
                        onOk={handleOk}
                        onCancel={handleCancel}
                    >
                        <p>Some contents...</p>
                        <p>Some contents...</p>
                        <p>Some contents...</p>
                    </Modal>
                </Content>
            </Layout>
        </>
    );
}


export default Settings;
