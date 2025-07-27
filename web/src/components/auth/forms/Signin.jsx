import { Card, Typography, Form, Input, Checkbox, Button } from 'antd';
const { Title } = Typography;
import { UserOutlined, LockOutlined } from '@ant-design/icons';
import * as auth from '../../../lib/auth';
import { useNavigate } from 'react-router';
import { useState } from 'react';
import { Alert } from 'antd';

const rememberedUsername = window.sessionStorage.getItem('username');
export default function () {
    const navigate = useNavigate();

    const [alert, setAlert] = useState('');
    const [loading, setLoading] = useState(false);

    const onSuccess = async values => {
        // TODO how to ensure csurf protection?
        setLoading(true); 
        const [err, data] = await auth.login(values.username, values.password);
        if (err) {
            setAlert('Something went wrong. Please try again later.');
            setLoading(false); 
            return;
        }

        if (data?.message) {
            setAlert(data.message);
            setLoading(false); 
            return;
        }
        
        setAlert('');
        // logged in successfully
        if (data.ID) {
            if (values.remember) {
                window.sessionStorage.setItem('username', values.username);
            }


            navigate('/');
        }

    };

    const onFail = errorInfo => {
        console.log('Failed:', errorInfo);
    };

    const LoginAlert = alert ? <Alert type="error" message={alert} showIcon={true} /> : '';

    return (
        <Card style={{ width: 500 }} >
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
                <Form.Item>
                    {LoginAlert}
                </ Form.Item>
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
                        loading={loading}
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
