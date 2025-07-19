import React, { useState } from 'react';
import { Layout, Menu, Button, Modal, theme, Breadcrumb, Space, Form, Input, Checkbox } from 'antd';
import { LaptopOutlined, NotificationOutlined, UserOutlined } from '@ant-design/icons';
const { Content, Sider } = Layout;

const onFinish = values => {
    console.log('Success:', values);
};
const onFinishFailed = errorInfo => {
    console.log('Failed:', errorInfo);
};

function Settings() {
    const [redditForm] = new Form.useForm();

    const [isModalOpen, setIsModalOpen] = useState(false);
    const [isModalLoading, setIsModalLoading] = useState(false);

    const showModal = () => {
        setIsModalOpen(true);
    };
    const handleOk = () => {
        setIsModalOpen(false);
    };
    const handleCancel = () => {
        setIsModalOpen(false);
    };

    const handleRedditConnection = () => {
        // TODO set loading state?
        console.log('connecting ...');
        
        redditForm.validateFields()
            .then((values) => {
                const username = values.username;
                const password = values.password;
                
                setIsModalLoading(true);
                setTimeout(function () {
                    console.log(username, password);
                    setIsModalLoading(false);
                }, 1500);
            }).catch((errorInfo) => {
                /*
                errorInfo: {
                    values: {
                        username: 'username',
                        password: 'password',
                    },
                    errorFields: [
                        { name: ['password'], errors: ['Please input your Password!'] },
                    ],
                    outOfDate: false,
                }
                */
            });
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
                            <i>
                            </i>
                            Reddit Connection
                        </label>
                        <Button type="primary" onClick={showModal}>Setup</Button>
                    </ Space>
                    <Modal
                        title="Connect to Reddit"
                        closable={{ 'aria-label': 'Custom Close Button' }}
                        open={isModalOpen}
                        onOk={handleOk}
                        onCancel={handleCancel}
                        loading={isModalLoading}
                        footer={[
                            <Button type="primary" onClick={handleRedditConnection}>
                                Connect
                            </Button>
                        ]}
                    >
                        <Form
                            name="basic"
                            initialValues={{ remember: true }}
                            form={redditForm}
                            onFinish={onFinish}
                            onFinishFailed={onFinishFailed}
                            autoComplete="off"
                        >
                            <Form.Item
                                label="Username"
                                name="username"
                                rules={[{ required: true, message: 'Please input your username!' }]}
                            >
                                <Input />
                            </Form.Item>

                            <Form.Item
                                label="Password"
                                name="password"
                                rules={[{ required: true, message: 'Please input your password!' }]}
                            >
                                <Input.Password />
                            </Form.Item>
                        </Form>
                    </Modal>
                </Content>
            </Layout>
        </>
    );
}


export default Settings;
