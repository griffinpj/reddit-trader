import { Card, Typography, Form, Input, Checkbox, Button } from 'antd';
const { Title } = Typography;
import { UserOutlined, LockOutlined } from '@ant-design/icons';
import * as auth from '../../../lib/auth';
import { useNavigate } from 'react-router';

const rememberedUsername = window.sessionStorage.getItem('username');


export default function () {
    const navigate = useNavigate();

    const onSuccess = async values => {
        // TODO how to ensure csurf protection?
        const data = await auth.login(values.username, values.password);

        // logged in successfully
        if (data.ID) {
            if (values.remember) {
                window.sessionStorage.setItem('username', values.username);
            }


            navigate('/');
        }

        // TODO show an alert for failed login data.message
    };

    const onFail = errorInfo => {
        console.log('Failed:', errorInfo);
    };

    return (
        <Card style={{ width: 500 }}>
            <div style={{ display: "flex", justifyContent: "center" }}>
                <Title level={2}>Log In</Title>
            </div>
            <Form
                name="login"
                initialValues={{ username: rememberedUsername, remember: true }}
                onFinish={onSuccess}
                onFinishFailed={onFail}
                autoComplete="off"
            >
                <Form.Item
                    name="username"
                    rules={[{ required: true, message: 'Please input your username!' }]}
                >
                    <Input placeholder="Username" prefix={<UserOutlined />} />
                </Form.Item>

                <Form.Item
                    style={{ paddingBottom: 0, marginBottom: 0 }}
                >
                    <Form.Item
                        noStyle
                        name="password"
                        rules={[{ required: true, message: 'Please input your password!' }]}
                    >
                        <Input.Password placeholder="Password" prefix={<LockOutlined />} />
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
                <a onClick={() => { navigate('/register'); }}> 
                    sign up
                </a>
            </Form>
        </Card>
    );
}
