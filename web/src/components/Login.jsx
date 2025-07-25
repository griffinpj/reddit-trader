import { Layout, Button, Checkbox, Form, Input, Flex, Typography, theme, Card } from 'antd';
import { UserOutlined, LockOutlined } from '@ant-design/icons';
const { Content, Header } = Layout;
const { Title } = Typography;

import * as auth from '../lib/auth';

const onFinishLogin = async values => {
    // TODO make ajax request to server with credentials
    // TODO how to ensure csurf protection?
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

const onFinishRegister = async values => {
};

const onFinishFailedRegister = errorInfo => {
    console.log('Failed:', errorInfo);
};

const onFinishFailedLogin = errorInfo => {
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
                        justify='center'
                        style={{
                            padding: 24
                        }}
                    >
                        <Card style={{ width: 500 }}>
                            <div style={{ display: "flex", justifyContent: "center" }}>
                                <Title level={2}>Log In</Title>
                            </div>
                            <Form
                                name="login"
                                initialValues={{ username: rememberedUsername, remember: true }}
                                onFinish={onFinishLogin}
                                onFinishFailed={onFinishFailedLogin}
                                autoComplete="off"
                            >
                                <Form.Item
                                    name="username"
                                    rules={[{ required: true, message: 'Please input your username!' }]}
                                >
                                    <Input prefix={<UserOutlined />} />
                                </Form.Item>

                                <Form.Item
                                    style={{ paddingBottom: 0, marginBottom: 0 }}
                                >
                                    <Form.Item
                                        noStyle
                                        name="password"
                                        rules={[{ required: true, message: 'Please input your password!' }]}
                                    >
                                        <Input.Password prefix={<LockOutlined />} />
                                    </Form.Item>
                                    <a
                                        style={{ float: "right" }}
                                        className="login-form-forgot"
                                        href=""
                                    >
                                        Forgot password
                                    </a>
                                </Form.Item>

                                <Form.Item name="remember" valuePropName="checked" label={null}>
                                    <Checkbox>Remember me</Checkbox>
                                </Form.Item>

                                <Form.Item label={null}>
                                    <Button
                                        type="primary"
                                        htmlType="submit"
                                        className="login-form-button"
                                        block
                                    >
                                        Log in
                                    </Button>
                                </Form.Item>
                                Don't have an account{" "}
                                <a href=""> 
                                    sign up
                                </a>
                            </Form>
                        </Card>
                    </Flex>
                </div>
            </Content>
        </Layout>
    );
};

export default Login;
