import { Card, Typography, Form, Input, Checkbox, Button, Divider } from 'antd';
const { Title } = Typography;
import { UserOutlined, LockOutlined, MailOutlined, CheckOutlined } from '@ant-design/icons';
import * as auth from '../../../lib/auth';
import { useNavigate } from 'react-router';


export default function () {
    
    const navigate = useNavigate();
    const onSuccess = async values => {
        // TODO how to ensure csurf protection?
        const [err, data] = await auth.register(values);
        if (err) {
            console.error(err);
            return;
        }

        // TODO show an alert for failed login data.message
        if (data?.message) {
            return;    
        }

        // logged in successfully
        navigate('/login');
    };

    const onFail = errorInfo => {
        console.log('Failed:', errorInfo);
    };

    return (
        <Card style={{ width: 500 }}>
            <div style={{ display: "flex", justifyContent: "center" }}>
                <Title level={2}>Register</Title>
            </div>
            <Form
                name="register"
                onFinish={onSuccess}
                onFinishFailed={onFail}
                autoComplete="off"
            >
                <Form.Item
                    name="first_name"
                    required={true}
                    rules={[{ required: true, message: 'Please input your first name!' }]}
                >
                    <Input placeholder="First Name" type="text" prefix={<UserOutlined />} />
                </Form.Item>
                <Form.Item
                    name="last_name"
                    required={true}
                    rules={[{ required: true, message: 'Please input your email!' }]}
                >
                    <Input placeholder="Last Name" type="text" prefix={<UserOutlined />} />
                </Form.Item>
                <Form.Item
                    name="email"
                    required={true}
                    rules={[{ required: true, message: 'Please input your email!' }]}
                >
                    <Input placeholder="Email" type="email" prefix={<MailOutlined />} />
                </Form.Item>

                <Form.Item
                    name="username"
                    required={true}
                    rules={[{ required: true, message: 'Please input your username!' }]}
                >
                    <Input placeholder="Username" prefix={<UserOutlined />} />
                </Form.Item>
                <Divider />
                <Form.Item>
                    <Form.Item
                        noStyle
                        name="password"
                        rules={[{ required: true, message: 'Please input your password!' }]}
                    >
                        <Input.Password placeholder="Password" prefix={<LockOutlined />} />
                    </Form.Item>
                </Form.Item>
                
                <Form.Item
                >
                    <Form.Item
                        noStyle
                        name="password-confirm"
                        rules={[{ required: true, message: 'Confirm your password!' }]}
                    >
                        <Input.Password placeholder="Confirm Password" prefix={<CheckOutlined />} />
                    </Form.Item>
                </Form.Item>

                <Form.Item label={null}>
                    <Button
                        type="primary"
                        htmlType="submit"
                        className="register-form-button"
                        block
                    >
                        Register
                    </Button>
                </Form.Item>
                Back to{' '}
                <a onClick={() => { navigate('/login'); }}> 
                    Login
                </a>
            </Form>
        </Card>
    );
}
