import { Layout, Button, Checkbox, Form, Input, Flex, Typography, theme } from 'antd';
const { Content, Header } = Layout;
const { Title } = Typography;

import * as auth from '../lib/auth';

const onFinish = async values => {
    // TODO make ajax request to server with credentials
    // TODO how to ensure csurf protection?
    console.log('Success:', values);
    const data = await auth.login(values.username, values.password);
    
    // logged in successfully
    if (data.ID) { 
        if (values.remember) {
            window.sessionStorage.setItem('username', values.username);
        }


        window.location = '/';
    }

    // TODO show an alert for failed login data.message
};

const onFinishFailed = errorInfo => {
    console.log('Failed:', errorInfo);
};

const rememberedUsername = window.sessionStorage.getItem('username');

function Login() {
    const {
        token: { colorBgContainer, borderRadiusLG },
    } = theme.useToken();

    return (
        <Layout>
            <Header style={{ display: 'flex', alignItems: 'center' }}>
                <div className="demo-logo" />
            </Header>
            <Content>
                <div
                    style={{
                        background: colorBgContainer,
                        minHeight: 280,
                        padding: 24,
                        borderRadius: borderRadiusLG
                    }}
                >
                    <Flex
                        vertical
                        align="center"
                        style={{
                            padding: 24
                        }}
                    >
                        <Title level={2}>Log In</Title>
                        <Form
                            name="login"
                            labelCol={{ span: 8 }}
                            wrapperCol={{ span: 16 }}
                            style={{ maxWidth: 600 }}
                            initialValues={{ remember: true }}
                            onFinish={onFinish}
                            onFinishFailed={onFinishFailed}
                            autoComplete="off"
                        >
                            <Form.Item
                                label="Username"
                                name="username"
                                initialValue={rememberedUsername}
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

                            <Form.Item name="remember" valuePropName="checked" label={null}>
                                <Checkbox>Remember me</Checkbox>
                            </Form.Item>

                            <Form.Item label={null}>
                                <Button type="primary" htmlType="submit">
                                    Submit
                                </Button>
                            </Form.Item>
                        </Form>
                    </Flex>
                </div>
            </Content>
        </Layout>
    );
};

export default Login;
